-- Rollback initial schema

DROP INDEX IF EXISTS idx_drafts_status;
DROP INDEX IF EXISTS idx_drafts_language;
DROP TABLE IF EXISTS drafts;

DROP INDEX IF EXISTS idx_puzzles_language_status;
DROP INDEX IF EXISTS idx_puzzles_status;
DROP INDEX IF EXISTS idx_puzzles_language_date;
DROP TABLE IF EXISTS puzzles;
