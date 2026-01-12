package store

import (
	"context"
	"sync"
	"time"

	"lesmotsdatche/internal/domain"
)

// MemoryStore is an in-memory store implementation for testing.
type MemoryStore struct {
	puzzles *MemoryPuzzleRepository
	drafts  *MemoryDraftRepository
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		puzzles: &MemoryPuzzleRepository{
			puzzles: make(map[string]*domain.Puzzle),
		},
		drafts: &MemoryDraftRepository{
			drafts: make(map[string]*Draft),
		},
	}
}

func (s *MemoryStore) Puzzles() PuzzleRepository { return s.puzzles }
func (s *MemoryStore) Drafts() DraftRepository   { return s.drafts }
func (s *MemoryStore) Migrate(ctx context.Context) error { return nil }
func (s *MemoryStore) Close() error { return nil }

// MemoryPuzzleRepository is an in-memory puzzle repository.
type MemoryPuzzleRepository struct {
	mu      sync.RWMutex
	puzzles map[string]*domain.Puzzle
}

func (r *MemoryPuzzleRepository) Store(ctx context.Context, p *domain.Puzzle) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Clone to prevent mutation
	clone := *p
	if clone.CreatedAt.IsZero() {
		clone.CreatedAt = time.Now()
	}
	r.puzzles[p.ID] = &clone
	return nil
}

func (r *MemoryPuzzleRepository) Get(ctx context.Context, id string) (*domain.Puzzle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.puzzles[id]
	if !ok {
		return nil, ErrNotFound
	}

	// Clone to prevent mutation
	clone := *p
	return &clone, nil
}

func (r *MemoryPuzzleRepository) GetByDate(ctx context.Context, language, date string) (*domain.Puzzle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.puzzles {
		if p.Language == language && p.Date == date {
			clone := *p
			return &clone, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MemoryPuzzleRepository) List(ctx context.Context, filter PuzzleFilter) ([]*PuzzleSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*PuzzleSummary
	for _, p := range r.puzzles {
		// Apply filters
		if filter.Language != "" && p.Language != filter.Language {
			continue
		}
		if filter.Status != "" && p.Status != filter.Status {
			continue
		}
		if filter.Difficulty > 0 && p.Difficulty != filter.Difficulty {
			continue
		}
		if filter.FromDate != "" && p.Date < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && p.Date > filter.ToDate {
			continue
		}

		result = append(result, &PuzzleSummary{
			ID:         p.ID,
			Date:       p.Date,
			Language:   p.Language,
			Title:      p.Title,
			Author:     p.Author,
			Difficulty: p.Difficulty,
			Status:     p.Status,
		})

		if filter.Limit > 0 && len(result) >= filter.Limit {
			break
		}
	}

	return result, nil
}

func (r *MemoryPuzzleRepository) UpdateStatus(ctx context.Context, id string, status domain.PuzzleStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.puzzles[id]
	if !ok {
		return ErrNotFound
	}

	p.Status = status
	if status == domain.StatusPublished && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	return nil
}

func (r *MemoryPuzzleRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.puzzles[id]; !ok {
		return ErrNotFound
	}
	delete(r.puzzles, id)
	return nil
}

// MemoryDraftRepository is an in-memory draft repository.
type MemoryDraftRepository struct {
	mu     sync.RWMutex
	drafts map[string]*Draft
}

func (r *MemoryDraftRepository) Store(ctx context.Context, d *Draft) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clone := *d
	if clone.CreatedAt.IsZero() {
		clone.CreatedAt = time.Now()
	}
	clone.UpdatedAt = time.Now()
	r.drafts[d.ID] = &clone
	return nil
}

func (r *MemoryDraftRepository) Get(ctx context.Context, id string) (*Draft, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.drafts[id]
	if !ok {
		return nil, ErrNotFound
	}
	clone := *d
	return &clone, nil
}

func (r *MemoryDraftRepository) List(ctx context.Context, language string) ([]*DraftSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*DraftSummary
	for _, d := range r.drafts {
		if language != "" && d.Language != language {
			continue
		}
		result = append(result, &DraftSummary{
			ID:        d.ID,
			Language:  d.Language,
			Status:    d.Status,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}
	return result, nil
}

func (r *MemoryDraftRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	d, ok := r.drafts[id]
	if !ok {
		return ErrNotFound
	}
	d.Status = status
	d.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryDraftRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.drafts[id]; !ok {
		return ErrNotFound
	}
	delete(r.drafts, id)
	return nil
}
