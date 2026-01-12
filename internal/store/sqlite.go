package store

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"lesmotsdatche/internal/domain"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("record not found")

// SQLiteStore implements Store using SQLite.
type SQLiteStore struct {
	db      *sql.DB
	puzzles *sqlitePuzzleRepo
	drafts  *sqliteDraftRepo
}

// NewSQLiteStore creates a new SQLite store.
// Use ":memory:" for in-memory database, or a file path for persistent storage.
func NewSQLiteStore(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better performance
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if !strings.Contains(dsn, ":memory:") {
		if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
		}
	}

	store := &SQLiteStore{db: db}
	store.puzzles = &sqlitePuzzleRepo{db: db}
	store.drafts = &sqliteDraftRepo{db: db}

	return store, nil
}

// Puzzles returns the puzzle repository.
func (s *SQLiteStore) Puzzles() PuzzleRepository {
	return s.puzzles
}

// Drafts returns the draft repository.
func (s *SQLiteStore) Drafts() DraftRepository {
	return s.drafts
}

// Migrate runs database migrations.
func (s *SQLiteStore) Migrate(ctx context.Context) error {
	upSQL, err := migrationsFS.ReadFile("migrations/001_initial.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration: %w", err)
	}

	_, err = s.db.ExecContext(ctx, string(upSQL))
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// sqlitePuzzleRepo implements PuzzleRepository for SQLite.
type sqlitePuzzleRepo struct {
	db *sql.DB
}

func (r *sqlitePuzzleRepo) Store(ctx context.Context, p *domain.Puzzle) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal puzzle: %w", err)
	}

	var publishedAt *time.Time
	if p.PublishedAt != nil {
		publishedAt = p.PublishedAt
	}

	// Use INSERT with ON CONFLICT DO UPDATE to handle updates by ID
	// but still fail on duplicate (language, date) for different IDs
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO puzzles (id, date, language, title, author, difficulty, status, payload, created_at, published_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			date = excluded.date,
			language = excluded.language,
			title = excluded.title,
			author = excluded.author,
			difficulty = excluded.difficulty,
			status = excluded.status,
			payload = excluded.payload,
			published_at = excluded.published_at
	`, p.ID, p.Date, p.Language, p.Title, p.Author, p.Difficulty, p.Status, payload, p.CreatedAt, publishedAt)

	if err != nil {
		return fmt.Errorf("failed to store puzzle: %w", err)
	}

	return nil
}

func (r *sqlitePuzzleRepo) Get(ctx context.Context, id string) (*domain.Puzzle, error) {
	var payload []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT payload FROM puzzles WHERE id = ?
	`, id).Scan(&payload)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get puzzle: %w", err)
	}

	var puzzle domain.Puzzle
	if err := json.Unmarshal(payload, &puzzle); err != nil {
		return nil, fmt.Errorf("failed to unmarshal puzzle: %w", err)
	}

	return &puzzle, nil
}

func (r *sqlitePuzzleRepo) GetByDate(ctx context.Context, language, date string) (*domain.Puzzle, error) {
	var payload []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT payload FROM puzzles WHERE language = ? AND date = ?
	`, language, date).Scan(&payload)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get puzzle by date: %w", err)
	}

	var puzzle domain.Puzzle
	if err := json.Unmarshal(payload, &puzzle); err != nil {
		return nil, fmt.Errorf("failed to unmarshal puzzle: %w", err)
	}

	return &puzzle, nil
}

func (r *sqlitePuzzleRepo) List(ctx context.Context, filter PuzzleFilter) ([]*PuzzleSummary, error) {
	query := `SELECT id, date, language, title, author, difficulty, status FROM puzzles WHERE 1=1`
	args := []interface{}{}

	if filter.Language != "" {
		query += " AND language = ?"
		args = append(args, filter.Language)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.FromDate != "" {
		query += " AND date >= ?"
		args = append(args, filter.FromDate)
	}
	if filter.ToDate != "" {
		query += " AND date <= ?"
		args = append(args, filter.ToDate)
	}
	if filter.Difficulty > 0 {
		query += " AND difficulty = ?"
		args = append(args, filter.Difficulty)
	}

	query += " ORDER BY date DESC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list puzzles: %w", err)
	}
	defer rows.Close()

	var puzzles []*PuzzleSummary
	for rows.Next() {
		var p PuzzleSummary
		if err := rows.Scan(&p.ID, &p.Date, &p.Language, &p.Title, &p.Author, &p.Difficulty, &p.Status); err != nil {
			return nil, fmt.Errorf("failed to scan puzzle: %w", err)
		}
		puzzles = append(puzzles, &p)
	}

	return puzzles, rows.Err()
}

func (r *sqlitePuzzleRepo) UpdateStatus(ctx context.Context, id string, status domain.PuzzleStatus) error {
	// First get the current puzzle to update its payload
	puzzle, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	// Update the status in the puzzle struct
	puzzle.Status = status
	if status == domain.StatusPublished && puzzle.PublishedAt == nil {
		now := time.Now().UTC()
		puzzle.PublishedAt = &now
	}

	// Re-marshal the payload with updated status
	payload, err := json.Marshal(puzzle)
	if err != nil {
		return fmt.Errorf("failed to marshal updated puzzle: %w", err)
	}

	result, err := r.db.ExecContext(ctx, `
		UPDATE puzzles SET status = ?, published_at = ?, payload = ? WHERE id = ?
	`, status, puzzle.PublishedAt, payload, id)

	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *sqlitePuzzleRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM puzzles WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete puzzle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// sqliteDraftRepo implements DraftRepository for SQLite.
type sqliteDraftRepo struct {
	db *sql.DB
}

func (r *sqliteDraftRepo) Store(ctx context.Context, d *Draft) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	d.UpdatedAt = time.Now().UTC()

	if d.Status == "" {
		d.Status = "draft"
	}

	payload, err := json.Marshal(d.Puzzle)
	if err != nil {
		return fmt.Errorf("failed to marshal puzzle: %w", err)
	}

	var report []byte
	if d.Report != nil {
		report, err = json.Marshal(d.Report)
		if err != nil {
			return fmt.Errorf("failed to marshal report: %w", err)
		}
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO drafts (id, language, payload, report, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			language = excluded.language,
			payload = excluded.payload,
			report = excluded.report,
			status = excluded.status,
			updated_at = excluded.updated_at
	`, d.ID, d.Language, payload, report, d.Status, d.CreatedAt, d.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to store draft: %w", err)
	}

	return nil
}

func (r *sqliteDraftRepo) Get(ctx context.Context, id string) (*Draft, error) {
	var d Draft
	var payload, report []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT id, language, payload, report, status, created_at, updated_at
		FROM drafts WHERE id = ?
	`, id).Scan(&d.ID, &d.Language, &payload, &report, &d.Status, &d.CreatedAt, &d.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	if err := json.Unmarshal(payload, &d.Puzzle); err != nil {
		return nil, fmt.Errorf("failed to unmarshal puzzle: %w", err)
	}

	if report != nil {
		d.Report = &domain.DraftReport{}
		if err := json.Unmarshal(report, d.Report); err != nil {
			return nil, fmt.Errorf("failed to unmarshal report: %w", err)
		}
	}

	return &d, nil
}

func (r *sqliteDraftRepo) List(ctx context.Context, language string) ([]*DraftSummary, error) {
	query := `SELECT id, language, status, created_at, updated_at FROM drafts`
	args := []interface{}{}

	if language != "" {
		query += " WHERE language = ?"
		args = append(args, language)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list drafts: %w", err)
	}
	defer rows.Close()

	var drafts []*DraftSummary
	for rows.Next() {
		var d DraftSummary
		if err := rows.Scan(&d.ID, &d.Language, &d.Status, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan draft: %w", err)
		}
		drafts = append(drafts, &d)
	}

	return drafts, rows.Err()
}

func (r *sqliteDraftRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE drafts SET status = ?, updated_at = ? WHERE id = ?
	`, status, time.Now().UTC(), id)

	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *sqliteDraftRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM drafts WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete draft: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
