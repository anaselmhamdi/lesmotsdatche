# Les Mots d'Atche

A modern French-first crossword puzzle application featuring LLM-assisted generation.

## Quick Start

```bash
# 1. Setup environment
cp .env.example .env
# Edit .env with your OPENAI_API_KEY

# 2. Run with Docker Compose
docker compose up --build

# Access:
# - API: http://localhost:8080
# - Web: http://localhost:3000
```

## Project Structure

```
cmd/
├── api/           # REST API server
└── generate/      # Puzzle generation CLI

internal/
├── domain/        # Core puzzle types
├── generator/     # Generation pipeline
│   ├── orchestrator.go  # Main coordinator
│   ├── llm/       # LLM client with validation
│   ├── fill/      # Constraint-based solver
│   ├── theme/     # Theme generation
│   ├── clue/      # Clue generation
│   ├── qa/        # Quality scoring
│   └── languagepack/  # FR/EN rules
├── api/           # HTTP handlers
├── store/         # SQLite persistence
└── validate/      # JSON schema validation

flutter_app/       # Mobile/web client
schemas/           # JSONSchema definitions
seed/              # Test data
testdata/          # Golden test fixtures
```

## Development

### Prerequisites
- Go 1.24+
- Flutter 3.0+
- OpenAI API key
- pre-commit (`pip install pre-commit && pre-commit install`)

### Commands

| Task | Command |
|------|---------|
| Build API | `go build -o ./api ./cmd/api` |
| Build Generator | `go build -o ./generator ./cmd/generate` |
| Run API | `go run ./cmd/api` |
| Run Tests | `go test ./...` |
| Generate Puzzle | `go run ./cmd/generate -lang fr -difficulty 3` |
| Flutter Dev | `cd flutter_app && flutter run` |

### Generator CLI Flags

```
-date        Target date (default: today)
-lang        Language: fr|en (default: fr)
-difficulty  1-5 (default: 3)
-output      Output file (default: stdout)
-api-key     OpenAI key (or use OPENAI_API_KEY env)
-model       Model name (default: gpt-4o)
-timeout     Generation timeout (default: 5m)
-max-attempts  Retry attempts (default: 3)
-verbose     Enable debug logging
```

### Before Committing / Creating PRs

Pre-commit hooks (`.pre-commit-config.yaml`) automatically run `go fmt`, `go vet`, `go mod tidy`, and file formatting checks.

Always run tests before committing:

```bash
# Run all Go tests
go test ./...

# Run Flutter tests (if flutter_app/ was modified)
cd flutter_app && flutter test

# Verify build succeeds
go build ./...
```

Update relevant documentation (README.md, CLAUDE.md, BRIEF.md) when changing:
- CLI flags or environment variables
- API endpoints or behavior
- Generation pipeline logic
- Language pack features

## Architecture

**Hybrid Generation Approach:**
- **LLM handles**: Themes, word candidates, clue writing, difficulty calibration
- **Deterministic code handles**: Grid filling, validation, QA scoring, safety checks

**Generation Pipeline:**
1. Theme Generation → LLM creates theme with keywords
2. Candidate Generation → LLM suggests words per slot length
3. Grid Filling → Constraint solver fills grid (backtracking algorithm)
4. Clue Generation → LLM writes clue variants
5. QA Scoring → Quality and safety evaluation

## API Reference

### Public Endpoints
- `GET /health` - Health check
- `GET /v1/puzzles/daily?language=fr` - Today's puzzle
- `GET /v1/puzzles?language=fr&from=&to=&difficulty=` - List puzzles
- `GET /v1/puzzles/{id}` - Get puzzle

### Admin Endpoints
- `POST /admin/v1/puzzles` - Store puzzle
- `PATCH /admin/v1/puzzles/{id}/status` - Update status
- `GET /admin/v1/puzzles` - List all puzzles

## Configuration

Environment variables (see `.env.example`):
- `OPENAI_API_KEY` - Required for generation
- `PORT` - Server port (default: `:8080`)
- `DATABASE_PATH` - SQLite file (default: `puzzles.db`)

## Internationalization

- **Current**: French (FR) fully implemented
- **Future**: English (EN) stub ready in `languagepack/en.go`
- Grid uses A-Z only; accents stripped but preserved in clue display
- Add new languages by implementing `LanguagePack` interface

## Documentation

- `BRIEF.md` - Full product specification
- `PRs.md` - Implementation breakdown (15 PRs)
- `CLAUDE.md` - AI assistant instructions
