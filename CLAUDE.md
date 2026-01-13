# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Les Mots d'Atche is a French-first crossword puzzle application with LLM-assisted generation. The system uses a hybrid approach: LLMs handle creative tasks (themes, clues) while deterministic code ensures grid correctness.

## Build & Run Commands

### Go Backend

```bash
# Build
go build -o ./api ./cmd/api
go build -o ./generator ./cmd/generate

# Run API server
go run ./cmd/api

# Run tests
go test ./...

# Run single test
go test -run TestAssignNumbers ./internal/domain

# Generate puzzle (requires OPENAI_API_KEY)
go run ./cmd/generate -lang fr -difficulty 3 -output puzzle.json -verbose
```

### Flutter App

```bash
cd flutter_app
flutter pub get
flutter test
flutter run

# After modifying freezed models
flutter pub run build_runner build
```

### Docker

```bash
# Full stack (API + seed + web) - production build
docker compose up --build

# API only
docker build -f Dockerfile.api -t lesmotsdatche-api .
```

### Development with Hot-Reload

For development with automatic reloading on code changes:

```bash
docker compose -f compose.yaml -f compose.dev.yaml up
```

- **Go**: Uses Air for hot-reload (rebuilds automatically on save)
- **Flutter**: Uses flutter web-server with hot-reload support

**NOTE**: Do not rebuild containers when making code changes in dev mode - hot-reload handles it automatically. Only rebuild if you change dependencies (go.mod, pubspec.yaml) or Dockerfiles.

## Architecture

### Generation Pipeline (internal/generator/)

The `Orchestrator` coordinates 7 sequential steps:
1. **Theme Generation** (LLM) → theme, keywords, seed words
2. **Template Creation** → 5x5 grid with symmetric blocks
3. **Slot Discovery** → identifies across/down entries
4. **Candidate Generation** (LLM) → word candidates per slot length
5. **Grid Filling** (deterministic) → constraint-based backtracking solver
6. **Clue Generation** (LLM) → clue variants for each answer
7. **QA Scoring** → quality evaluation and safety checks

Key design: LLMs handle creativity, deterministic solver guarantees grid correctness.

### Package Structure

- `internal/domain/` - Puzzle types, normalization, numbering
- `internal/generator/llm/` - LLM client with JSON schema validation and retry
- `internal/generator/fill/` - Constraint-based grid solver (NOT LLM)
- `internal/generator/theme/` - Theme and candidate generation
- `internal/generator/clue/` - Clue variant generation
- `internal/generator/qa/` - Quality scoring and safety filters
- `internal/generator/languagepack/` - Language-specific rules (FR implemented, EN stub)
- `internal/api/` - REST handlers and middleware
- `internal/store/` - SQLite repository layer

### API Endpoints

**Public:**
- `GET /health` - Health check
- `GET /v1/puzzles/daily?language=fr` - Today's puzzle
- `GET /v1/puzzles/{id}` - Get puzzle by ID

**Admin:**
- `POST /admin/v1/puzzles` - Store puzzle
- `PATCH /admin/v1/puzzles/{id}/status` - Update status

## Environment Variables

```
OPENAI_API_KEY=sk-...   # Required for generation
PORT=:8080              # API server port
DATABASE_PATH=puzzles.db # SQLite file path
```

## Key Patterns

- **Language packs**: Extensible via `languagepack.Register()` - FR implemented, EN stub ready
- **ValidatingClient**: Wraps LLM calls with JSON schema validation and auto-retry
- **A-Z only grids**: Accents stripped via `Normalize()`, display text preserved in clues
- **QA gating**: Puzzles must pass quality threshold before publishing

## Test Data

Golden fixtures in `testdata/`:
- `valid_7x7.json` - Complete valid puzzle
- `small_5x5_*.json` - Grid variations
- `invalid_*.json` - Error case fixtures

## Before Committing / Creating PRs

Pre-commit hooks are configured (`.pre-commit-config.yaml`) and will automatically run:
- `go fmt`, `go vet`, `go mod tidy`
- Trailing whitespace and EOF fixes
- YAML/JSON validation

Additionally, always run tests before committing:

```bash
# Run all Go tests
go test ./...

# Run Flutter tests (if flutter_app/ was modified)
cd flutter_app && flutter test

# Verify build succeeds
go build ./...
```

Update documentation when:
- Adding new CLI flags or environment variables
- Changing API endpoints
- Modifying the generation pipeline
- Adding new language pack features

## Common Tasks

**Seed the API with test data:**
```bash
curl -X POST http://localhost:8080/admin/v1/puzzles \
  -H "Content-Type: application/json" \
  -d @seed/puzzle.json
```

**Fix go.mod issues:**
```bash
go mod tidy
```
