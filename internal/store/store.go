// Package store provides database storage for puzzles and drafts.
package store

import (
	"context"
	"time"

	"lesmotsdatche/internal/domain"
)

// PuzzleFilter contains criteria for listing puzzles.
type PuzzleFilter struct {
	Language   string
	Status     domain.PuzzleStatus
	FromDate   string // YYYY-MM-DD
	ToDate     string // YYYY-MM-DD
	Tag        string
	Difficulty int
	Limit      int
	Offset     int
}

// PuzzleSummary contains summary info for puzzle listings.
type PuzzleSummary struct {
	ID         string       `json:"id"`
	Date       string       `json:"date"`
	Language   string       `json:"language"`
	Title      string       `json:"title"`
	Author     string       `json:"author"`
	Difficulty int          `json:"difficulty"`
	Status     domain.PuzzleStatus `json:"status"`
}

// DraftSummary contains summary info for draft listings.
type DraftSummary struct {
	ID        string    `json:"id"`
	Language  string    `json:"language"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Draft represents a puzzle draft with its QA report.
type Draft struct {
	ID        string              `json:"id"`
	Language  string              `json:"language"`
	Puzzle    domain.Puzzle       `json:"puzzle"`
	Report    *domain.DraftReport `json:"report,omitempty"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// PuzzleRepository defines the interface for puzzle storage operations.
type PuzzleRepository interface {
	// Store saves a puzzle to the database.
	Store(ctx context.Context, p *domain.Puzzle) error

	// Get retrieves a puzzle by ID.
	Get(ctx context.Context, id string) (*domain.Puzzle, error)

	// GetByDate retrieves a puzzle by language and date.
	GetByDate(ctx context.Context, language, date string) (*domain.Puzzle, error)

	// List returns puzzles matching the filter criteria.
	List(ctx context.Context, filter PuzzleFilter) ([]*PuzzleSummary, error)

	// UpdateStatus changes the status of a puzzle.
	UpdateStatus(ctx context.Context, id string, status domain.PuzzleStatus) error

	// Delete removes a puzzle by ID.
	Delete(ctx context.Context, id string) error
}

// DraftRepository defines the interface for draft storage operations.
type DraftRepository interface {
	// Store saves a draft to the database.
	Store(ctx context.Context, d *Draft) error

	// Get retrieves a draft by ID.
	Get(ctx context.Context, id string) (*Draft, error)

	// List returns drafts for a language.
	List(ctx context.Context, language string) ([]*DraftSummary, error)

	// UpdateStatus changes the status of a draft.
	UpdateStatus(ctx context.Context, id string, status string) error

	// Delete removes a draft by ID.
	Delete(ctx context.Context, id string) error
}

// Store combines all repository interfaces.
type Store interface {
	Puzzles() PuzzleRepository
	Drafts() DraftRepository

	// Migrate runs database migrations.
	Migrate(ctx context.Context) error

	// Close closes the database connection.
	Close() error
}
