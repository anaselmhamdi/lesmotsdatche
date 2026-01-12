-- Initial schema for crossword puzzle storage

-- Published puzzles table
CREATE TABLE IF NOT EXISTS puzzles (
    id TEXT PRIMARY KEY,
    date TEXT NOT NULL,
    language TEXT NOT NULL CHECK (language IN ('fr', 'en')),
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    difficulty INTEGER NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
    status TEXT NOT NULL CHECK (status IN ('draft', 'published', 'archived')),
    payload JSON NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    UNIQUE(language, date)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_puzzles_language_date ON puzzles(language, date);
CREATE INDEX IF NOT EXISTS idx_puzzles_status ON puzzles(status);
CREATE INDEX IF NOT EXISTS idx_puzzles_language_status ON puzzles(language, status);

-- Draft puzzles table
CREATE TABLE IF NOT EXISTS drafts (
    id TEXT PRIMARY KEY,
    language TEXT NOT NULL CHECK (language IN ('fr', 'en')),
    payload JSON NOT NULL,
    report JSON,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'rejected', 'published')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_drafts_language ON drafts(language);
CREATE INDEX IF NOT EXISTS idx_drafts_status ON drafts(status);
