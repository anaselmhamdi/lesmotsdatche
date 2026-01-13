# Les Mots d'Atche - Architecture Documentation

This document provides comprehensive technical documentation for Les Mots d'Atche, a French-first crossword puzzle application with LLM-assisted generation.

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [System Architecture](#2-system-architecture)
3. [Generation Pipeline](#3-generation-pipeline)
4. [LLM Integration](#4-llm-integration)
5. [Grid Solver](#5-grid-solver)
6. [QA System](#6-qa-system)
7. [REST API](#7-rest-api)
8. [Database Layer](#8-database-layer)
9. [Language Packs](#9-language-packs)
10. [Flutter App](#10-flutter-app)
11. [Deployment](#11-deployment)
12. [Key Files Reference](#12-key-files-reference)

---

## 1. Project Overview

### 1.1 What Is Les Mots d'Atche?

Les Mots d'Atche is a French-first crossword puzzle application that combines AI-assisted generation with a mobile/web interface for puzzle solving. The system produces high-quality daily puzzles through a hybrid approach:

- **LLMs handle creativity**: Theme selection, word candidates, clue writing
- **Deterministic code ensures correctness**: Grid filling, validation, quality scoring

### 1.2 Design Philosophy

The core insight is that LLMs excel at creative tasks but can produce inconsistent outputs. By using LLMs only where creativity matters and relying on deterministic algorithms for correctness, the system achieves both quality and reliability.

```
┌─────────────────────────────────────────────────────────────────┐
│                    HYBRID ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   LLM Tasks (Creative)          Deterministic Tasks (Reliable) │
│   ═══════════════════           ════════════════════════════   │
│   • Theme generation            • Grid template creation        │
│   • Word candidates             • Constraint-based filling      │
│   • Clue writing                • Slot discovery                │
│                                 • QA scoring                    │
│                                 • Normalization (A-Z)           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 1.3 Technology Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.24 |
| LLM Provider | OpenAI (GPT-4o) |
| Database | SQLite with WAL mode |
| Mobile App | Flutter (Dart) |
| State Management | Riverpod |
| Containerization | Docker Compose |

---

## 2. System Architecture

### 2.1 High-Level Component Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                      │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│    ┌─────────────┐     ┌─────────────┐     ┌─────────────┐              │
│    │ Flutter iOS │     │Flutter Andrd│     │  Flutter Web│              │
│    └──────┬──────┘     └──────┬──────┘     └──────┬──────┘              │
│           │                   │                   │                      │
│           └───────────────────┼───────────────────┘                      │
│                               │                                          │
│                               ▼                                          │
│                     ┌─────────────────┐                                  │
│                     │    REST API     │                                  │
│                     │   (Go/HTTP)     │                                  │
│                     └────────┬────────┘                                  │
│                              │                                           │
├──────────────────────────────┼───────────────────────────────────────────┤
│                              │           BACKEND                         │
│              ┌───────────────┼───────────────┐                          │
│              │               │               │                          │
│              ▼               ▼               ▼                          │
│    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                   │
│    │   Handlers  │  │    Store    │  │Orchestrator │                   │
│    │  (routes)   │  │  (sqlite)   │  │ (generator) │                   │
│    └─────────────┘  └──────┬──────┘  └──────┬──────┘                   │
│                            │                │                           │
│                            ▼                │                           │
│                     ┌───────────┐           │                           │
│                     │  SQLite   │           │                           │
│                     │    DB     │           │                           │
│                     └───────────┘           │                           │
│                                             │                           │
│         ┌───────────────────────────────────┼─────────────────┐        │
│         │                                   │                 │        │
│         ▼                                   ▼                 ▼        │
│  ┌─────────────┐                    ┌─────────────┐   ┌─────────────┐  │
│  │   Theme     │                    │    Fill     │   │    Clue     │  │
│  │ Generator   │                    │   Solver    │   │  Generator  │  │
│  └──────┬──────┘                    └─────────────┘   └──────┬──────┘  │
│         │                                                    │         │
│         │              ┌─────────────┐                       │         │
│         └──────────────►   OpenAI    ◄───────────────────────┘         │
│                        │    API      │                                  │
│                        └─────────────┘                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Package Structure

```
lesmotsdatche/
├── cmd/
│   ├── api/              # HTTP server binary
│   │   └── main.go
│   └── generate/         # CLI generation tool
│       └── main.go
├── internal/
│   ├── domain/           # Core types and business logic
│   │   ├── types.go      # Puzzle, Cell, Clue structs
│   │   ├── normalize.go  # Text normalization (→ A-Z)
│   │   └── numbering.go  # Clue numbering algorithm
│   ├── generator/        # Puzzle generation pipeline
│   │   ├── orchestrator.go
│   │   ├── llm/          # LLM client wrapper
│   │   ├── theme/        # Theme & candidate generation
│   │   ├── clue/         # Clue generation
│   │   ├── fill/         # Deterministic solver
│   │   ├── qa/           # Quality scoring
│   │   └── languagepack/ # Language-specific rules
│   ├── api/              # HTTP handlers and routing
│   │   ├── routes.go
│   │   ├── handlers.go
│   │   ├── admin.go
│   │   └── middleware.go
│   └── store/            # Database layer
│       ├── store.go      # Repository interfaces
│       ├── sqlite.go     # SQLite implementation
│       └── migrations/
├── flutter_app/          # Mobile/web client
├── schemas/              # JSON schemas for validation
├── seed/                 # Test data for seeding
└── testdata/             # Test fixtures
```

### 2.3 Data Flow Overview

```
Generation Flow:
────────────────
User Request → Orchestrator → [Theme→Template→Slots→Candidates→Fill→Clues→QA] → Puzzle

API Flow:
─────────
Flutter App → GET /v1/puzzles/daily → Handler → Store → SQLite → Response
```

---

## 3. Generation Pipeline

The `Orchestrator` coordinates 7 sequential steps to generate a complete puzzle. This is the heart of the system.

### 3.1 Pipeline Overview

```
┌──────────────────────────────────────────────────────────────────────────┐
│                        GENERATION PIPELINE                               │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Step 1: Theme Generation (LLM)                                          │
│  ═══════════════════════════════                                         │
│  Input:  Language, difficulty, date                                      │
│  Output: Theme { title, keywords, seed_words }                           │
│  File:   internal/generator/theme/theme.go                               │
│                                                                          │
│                              ▼                                           │
│                                                                          │
│  Step 2: Template Creation (Deterministic)                               │
│  ═════════════════════════════════════════                               │
│  Input:  Grid size (13x13 for French)                                    │
│  Output: [][]Cell with symmetric block pattern                           │
│  File:   internal/generator/orchestrator.go                              │
│                                                                          │
│                              ▼                                           │
│                                                                          │
│  Step 3: Slot Discovery (Deterministic)                                  │
│  ══════════════════════════════════════                                  │
│  Input:  Template grid                                                   │
│  Output: []Slot (across and down word positions)                         │
│  File:   internal/generator/fill/slot.go                                 │
│                                                                          │
│                              ▼                                           │
│                                                                          │
│  Step 4: Candidate Generation (LLM)                                      │
│  ══════════════════════════════════                                      │
│  Input:  Theme, slot lengths, language pack                              │
│  Output: MemoryLexicon with themed word candidates                       │
│  File:   internal/generator/theme/candidates.go                          │
│                                                                          │
│                              ▼                                           │
│                                                                          │
│  Step 5: Grid Filling (Deterministic)                                    │
│  ════════════════════════════════════                                    │
│  Input:  Template, lexicon, slots                                        │
│  Output: Filled grid with all letters                                    │
│  File:   internal/generator/fill/solver.go                               │
│                                                                          │
│                              ▼                                           │
│                                                                          │
│  Step 6: Clue Generation (LLM)                                           │
│  ═════════════════════════════                                           │
│  Input:  Filled grid, theme, answers                                     │
│  Output: Clue variants for each answer                                   │
│  File:   internal/generator/clue/clue.go                                 │
│                                                                          │
│                              ▼                                           │
│                                                                          │
│  Step 7: QA Scoring (Deterministic)                                      │
│  ══════════════════════════════════                                      │
│  Input:  Complete puzzle                                                 │
│  Output: Score (0.0-1.0) + flags                                         │
│  File:   internal/generator/qa/scorer.go                                 │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Orchestrator Configuration

```go
// internal/generator/orchestrator.go

type Config struct {
    MaxAttempts      int           // Default: 3
    Timeout          time.Duration // Default: 5 minutes
    TargetDifficulty int           // 1-5
    MinQAScore       float64       // Default: 0.5
    GridSize         [2]int        // Default: [13, 13] for French
}

func DefaultConfig() Config {
    return Config{
        MaxAttempts:      3,
        Timeout:          5 * time.Minute,
        TargetDifficulty: 3,
        MinQAScore:       0.5,
        GridSize:         [2]int{13, 13},
    }
}
```

### 3.3 Multi-Attempt Logic

The orchestrator makes multiple attempts to generate an acceptable puzzle:

```go
func (o *Orchestrator) Generate(ctx context.Context) (*GenerateResult, error) {
    var lastErr error

    for attempt := 1; attempt <= o.config.MaxAttempts; attempt++ {
        result, err := o.generateOnce(ctx)
        if err != nil {
            lastErr = err
            continue
        }

        if result.Score.IsAcceptable(o.config.MinQAScore) {
            return result, nil
        }

        lastErr = fmt.Errorf("QA score %.2f below threshold", result.Score.Overall)
    }

    return nil, fmt.Errorf("failed after %d attempts: %w", o.config.MaxAttempts, lastErr)
}
```

### 3.4 Step Details

#### Step 1: Theme Generation

The theme generator calls the LLM to produce a coherent theme with keywords and seed words.

```go
type Theme struct {
    Title       string   // e.g., "Le Cinéma Français"
    Description string   // Brief description
    Keywords    []string // Related terms for clue writing
    SeedWords   []string // Starting words to include
    Difficulty  int      // 1-5
}
```

#### Step 2: Template Creation

Creates a symmetric block pattern typical of French crosswords:

```
┌───┬───┬───┬───┬───┬───┬───┐
│   │   │   │ ■ │   │   │   │
├───┼───┼───┼───┼───┼───┼───┤
│   │   │   │   │   │   │   │
├───┼───┼───┼───┼───┼───┼───┤
│   │   │ ■ │   │ ■ │   │   │
├───┼───┼───┼───┼───┼───┼───┤
│ ■ │   │   │   │   │   │ ■ │
├───┼───┼───┼───┼───┼───┼───┤
│   │   │ ■ │   │ ■ │   │   │
├───┼───┼───┼───┼───┼───┼───┤
│   │   │   │   │   │   │   │
├───┼───┼───┼───┼───┼───┼───┤
│   │   │   │ ■ │   │   │   │
└───┴───┴───┴───┴───┴───┴───┘
```

#### Step 3: Slot Discovery

Identifies all word positions (across and down) in the template:

```go
type Slot struct {
    Direction Direction // Across or Down
    Start     Position  // Row, Col
    Length    int
    Cells     []Position
}

func DiscoverSlots(template [][]Cell) []Slot
```

#### Step 4: Candidate Generation

Generates themed word candidates grouped by length:

```go
type CandidateGenerator struct {
    client   *llm.ValidatingClient
    langPack languagepack.LanguagePack
}

// Returns candidates for each required word length
func (g *CandidateGenerator) GenerateCandidates(ctx context.Context, theme Theme, lengths []int) (*MemoryLexicon, error)
```

#### Step 5: Grid Filling

The deterministic solver fills the grid using constraint-based backtracking (see [Section 5](#5-grid-solver)).

#### Step 6: Clue Generation

Generates multiple clue variants per answer with different styles:

```go
type ClueCandidate struct {
    Prompt     string // The clue text
    Style      string // definition, wordplay, cultural
    Difficulty int    // 1-5
}
```

#### Step 7: QA Scoring

Evaluates puzzle quality across multiple dimensions (see [Section 6](#6-qa-system)).

---

## 4. LLM Integration

### 4.1 ValidatingClient Pattern

The `ValidatingClient` wraps the underlying LLM client to provide:
- JSON schema validation of responses
- Automatic retry with repair prompts
- Trace recording for debugging
- Secret redaction in logs

```
┌─────────────────────────────────────────────────────────────────┐
│                    ValidatingClient                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   Request    │ →  │  LLM Call    │ →  │   Validate   │      │
│  │              │    │              │    │   Response   │      │
│  └──────────────┘    └──────────────┘    └──────┬───────┘      │
│                                                  │              │
│                            ┌─────────────────────┴───────┐     │
│                            │                             │     │
│                            ▼                             ▼     │
│                     ┌──────────┐                 ┌──────────┐  │
│                     │  Valid   │                 │ Invalid  │  │
│                     │  Return  │                 │  Retry   │  │
│                     └──────────┘                 └────┬─────┘  │
│                                                       │        │
│                                                       ▼        │
│                                              ┌──────────────┐  │
│                                              │Repair Prompt │  │
│                                              │ + Re-submit  │  │
│                                              └──────────────┘  │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Configuration

```go
// internal/generator/llm/client.go

type Config struct {
    MaxRetries    int     // Default: 3
    DefaultTemp   float64 // Default: 0.7
    DefaultTokens int     // Default: 2048
    RedactSecrets bool    // Default: true
    RepairPrompt  string  // Template for retry attempts
}
```

### 4.3 Schema Validation

Responses are validated against JSON schemas before being accepted:

```go
func (c *ValidatingClient) CompleteWithValidation(
    ctx context.Context,
    req Request,
    target interface{},
) error {
    for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
        resp, err := c.client.Complete(ctx, req)
        if err != nil {
            return err
        }

        // Extract JSON from markdown code blocks
        content := extractJSON(resp.Content)

        // Validate against schema
        if err := req.Schema.Validate(content); err != nil {
            // Retry with repair prompt
            req.Prompt = fmt.Sprintf(c.config.RepairPrompt, err.Error())
            continue
        }

        return json.Unmarshal([]byte(content), target)
    }

    return ErrMaxRetries
}
```

### 4.4 Trace Recording

All LLM interactions are recorded for debugging:

```go
type Trace struct {
    Timestamp time.Time
    Request   Request
    Response  Response
    Duration  time.Duration
    Error     error
}

// Secrets are automatically redacted from traces
func (c *ValidatingClient) Traces() []Trace
```

### 4.5 Prompt Templates

Each language pack provides localized prompt templates:

```go
// internal/generator/languagepack/fr.go

type PromptTemplates struct {
    ThemeGeneration string // "Génère un thème de mots croisés..."
    SlotCandidates  string // "Propose des mots de {length} lettres..."
    ClueGeneration  string // "Écris des définitions pour..."
    ClueStyle       string // "Style: définition, jeu de mots, culturel"
}
```

---

## 5. Grid Solver

The fill solver is a constraint-based backtracking algorithm that guarantees a valid solution.

### 5.1 Algorithm Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    BACKTRACKING SOLVER                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. Select most-constrained slot (fewest valid candidates)      │
│                                                                 │
│  2. For each candidate word (in scored order):                  │
│     a. Place word in slot                                       │
│     b. Propagate constraints to crossing slots                  │
│     c. If valid: recurse to next slot                          │
│     d. If invalid or recursion fails: backtrack                │
│                                                                 │
│  3. If all slots filled: return success                        │
│                                                                 │
│  4. If backtrack limit reached: return failure                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 Solver Structure

```go
// internal/generator/fill/solver.go

type Solver struct {
    lexicon        Lexicon // Word database
    scorer         Scorer  // Candidate ranking
    rng            *rand.Rand
    maxBacktrack   int     // Default: 50,000
    backtrackCount int
}

func NewSolver(lexicon Lexicon, opts ...Option) *Solver

func (s *Solver) Solve(template [][]Cell) ([][]Cell, error)
```

### 5.3 Most-Constrained-First Heuristic

The solver prioritizes slots with fewer valid candidates:

```go
func (s *Solver) selectNextSlot(slots []Slot, grid [][]Cell) *Slot {
    var best *Slot
    minCandidates := math.MaxInt

    for _, slot := range slots {
        if slot.IsFilled(grid) {
            continue
        }

        pattern := slot.GetPattern(grid)
        candidates := s.lexicon.Match(pattern)

        if len(candidates) < minCandidates {
            best = &slot
            minCandidates = len(candidates)
        }
    }

    return best
}
```

### 5.4 Lexicon Interface

```go
type Lexicon interface {
    // Match returns words matching pattern (dots are wildcards)
    // e.g., "C.T" matches "CAT", "COT", "CUT"
    Match(pattern string) []string

    // Contains checks if word exists
    Contains(word string) bool

    // ByLength returns all words of given length
    ByLength(length int) []string
}

type MemoryLexicon struct {
    words    map[string]WordEntry
    byLength map[int][]string
}
```

### 5.5 Candidate Scoring

Words are scored by theme relevance and frequency:

```go
type WordEntry struct {
    Word      string
    Frequency float64 // Usage frequency
    ThemeScore float64 // Relevance to current theme
}

func (s *Solver) scoreCandidates(candidates []string, slot Slot) []string {
    // Sort by combined score
    // Add randomization within tiers for variety
}
```

---

## 6. QA System

The QA system evaluates puzzle quality and enforces safety standards.

### 6.1 Scoring Components

```
┌─────────────────────────────────────────────────────────────────┐
│                      QA SCORING                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Component          Weight    Description                       │
│  ─────────────────────────────────────────────────────────────  │
│  Fill Score          25%      Word quality, backtrack penalty   │
│  Clue Score          30%      Variety, length, completeness     │
│  Freshness Score     20%      Uniqueness vs recent puzzles      │
│  Structure Score     25%      Symmetry, block density           │
│                                                                 │
│  Overall = Σ(component × weight)                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 6.2 Score Structure

```go
// internal/generator/qa/scorer.go

type Score struct {
    Overall    float64            // 0.0-1.0 aggregate score
    Components map[string]float64 // Individual component scores
    Flags      []Flag             // Warnings and errors
}

type Flag struct {
    Level   FlagLevel // info, warning, error
    Code    string    // Unique identifier
    Message string    // Human readable description
}

type FlagLevel string

const (
    FlagInfo    FlagLevel = "info"
    FlagWarning FlagLevel = "warning"
    FlagError   FlagLevel = "error"
)
```

### 6.3 Acceptance Criteria

```go
func (s Score) IsAcceptable(threshold float64) bool {
    // Must meet minimum score
    if s.Overall < threshold {
        return false
    }

    // Must have no error-level flags
    for _, flag := range s.Flags {
        if flag.Level == FlagError {
            return false
        }
    }

    return true
}
```

### 6.4 Safety Filters

The QA system checks for inappropriate content:

```go
func (scorer *Scorer) checkSafety(puzzle *Puzzle) []Flag {
    var flags []Flag

    for _, clue := range puzzle.AllClues() {
        // Check answer against taboo list
        if scorer.langPack.IsTaboo(clue.Answer) {
            flags = append(flags, Flag{
                Level:   FlagError,
                Code:    "TABOO_WORD",
                Message: fmt.Sprintf("Taboo word in answer: %s", clue.Answer),
            })
        }

        // Check clue text for inappropriate content
        // ...
    }

    return flags
}
```

---

## 7. REST API

### 7.1 Endpoint Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        REST API                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  PUBLIC ENDPOINTS                                               │
│  ════════════════                                               │
│  GET  /health                    Health check                   │
│  GET  /v1/puzzles/daily          Today's puzzle (by language)   │
│  GET  /v1/puzzles/{id}           Get puzzle by ID               │
│  GET  /v1/puzzles                List puzzles (with filters)    │
│                                                                 │
│  ADMIN ENDPOINTS                                                │
│  ═══════════════                                                │
│  POST   /admin/v1/puzzles         Store new puzzle              │
│  PATCH  /admin/v1/puzzles/{id}/status  Update puzzle status     │
│  GET    /admin/v1/puzzles         List all (inc. drafts)        │
│  GET    /admin/v1/puzzles/{id}    Get puzzle (admin view)       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 7.2 Request/Response Examples

#### Get Daily Puzzle

```http
GET /v1/puzzles/daily?language=fr HTTP/1.1
Host: api.lesmotsdatche.com
Accept: application/json
```

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "date": "2024-01-15",
  "language": "fr",
  "title": "Le Cinéma Français",
  "author": "Les Mots d'Atche",
  "difficulty": 3,
  "status": "published",
  "grid": [
    [{"type": "letter", "solution": "C", "number": 1}, ...],
    ...
  ],
  "clues": {
    "across": [
      {"number": 1, "prompt": "Art du grand écran", "answer": "CINEMA", ...}
    ],
    "down": [...]
  }
}
```

#### Store Puzzle (Admin)

```http
POST /admin/v1/puzzles HTTP/1.1
Host: api.lesmotsdatche.com
Content-Type: application/json

{
  "date": "2024-01-16",
  "language": "fr",
  "title": "...",
  ...
}
```

### 7.3 Middleware Stack

```go
// internal/api/routes.go

func NewRouter(cfg Config) http.Handler {
    r := chi.NewRouter()

    // Middleware (applied in order)
    r.Use(middleware.Recover)    // Panic recovery
    r.Use(middleware.Logger)     // Request logging
    r.Use(middleware.Gzip)       // Response compression
    r.Use(middleware.CORS)       // Cross-origin support

    // Routes...
    return r
}
```

### 7.4 Caching

Responses include ETag and cache headers:

```go
func (h *Handler) GetDaily(w http.ResponseWriter, r *http.Request) {
    puzzle, err := h.store.Puzzles().GetByDate(language, today)

    // Generate ETag from puzzle ID and updated timestamp
    etag := fmt.Sprintf(`"%s-%d"`, puzzle.ID, puzzle.UpdatedAt.Unix())

    w.Header().Set("ETag", etag)
    w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes

    // Check If-None-Match for conditional request
    if r.Header.Get("If-None-Match") == etag {
        w.WriteHeader(http.StatusNotModified)
        return
    }

    json.NewEncoder(w).Encode(puzzle)
}
```

---

## 8. Database Layer

### 8.1 Repository Pattern

```go
// internal/store/store.go

type Store interface {
    Puzzles() PuzzleRepository
    Drafts()  DraftRepository
    Migrate(ctx context.Context) error
    Close() error
}

type PuzzleRepository interface {
    Store(ctx context.Context, p *Puzzle) error
    Get(ctx context.Context, id string) (*Puzzle, error)
    GetByDate(ctx context.Context, lang, date string) (*Puzzle, error)
    List(ctx context.Context, filter Filter) ([]*PuzzleSummary, error)
    UpdateStatus(ctx context.Context, id string, status Status) error
    Delete(ctx context.Context, id string) error
}
```

### 8.2 SQLite Schema

```sql
-- Puzzles table
CREATE TABLE puzzles (
    id TEXT PRIMARY KEY,
    date TEXT NOT NULL,
    language TEXT NOT NULL,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    difficulty INTEGER,
    status TEXT DEFAULT 'draft',
    payload JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    UNIQUE(language, date)
);

-- Indexes
CREATE INDEX idx_puzzles_language_date ON puzzles(language, date);
CREATE INDEX idx_puzzles_status ON puzzles(status);
CREATE INDEX idx_puzzles_language_status ON puzzles(language, status);

-- Drafts table (for in-progress generation)
CREATE TABLE drafts (
    id TEXT PRIMARY KEY,
    language TEXT NOT NULL,
    payload JSON NOT NULL,
    report JSON,
    status TEXT DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 8.3 JSON Payload Storage

Complex puzzle data is stored as a JSON blob with indexed columns for common queries:

```go
func (r *sqlitePuzzleRepo) Store(ctx context.Context, p *Puzzle) error {
    payload, _ := json.Marshal(p)

    _, err := r.db.ExecContext(ctx, `
        INSERT INTO puzzles (id, date, language, title, author, difficulty, status, payload)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            payload = excluded.payload,
            status = excluded.status
    `, p.ID, p.Date, p.Language, p.Title, p.Author, p.Difficulty, p.Status, payload)

    return err
}
```

### 8.4 WAL Mode

SQLite is configured with Write-Ahead Logging for better concurrency:

```go
func NewSQLiteStore(path string) (*SQLiteStore, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, err
    }

    // Enable WAL mode (except for in-memory DBs)
    if path != ":memory:" {
        db.Exec("PRAGMA journal_mode = WAL")
    }

    db.Exec("PRAGMA foreign_keys = ON")

    return &SQLiteStore{db: db}, nil
}
```

---

## 9. Language Packs

### 9.1 Interface

```go
// internal/generator/languagepack/pack.go

type LanguagePack interface {
    Code() string              // "fr", "en"
    Name() string              // "Français", "English"
    Normalize(text string) string
    IsTaboo(word string) bool
    TabooList() []string
    IsConfigured() bool
    Prompts() PromptTemplates
}

type PromptTemplates struct {
    ThemeGeneration string
    SlotCandidates  string
    ClueGeneration  string
    ClueStyle       string
}
```

### 9.2 Registry

```go
type Registry struct {
    packs map[string]LanguagePack
}

func DefaultRegistry() *Registry {
    r := &Registry{packs: make(map[string]LanguagePack)}
    r.Register(NewFrenchPack())
    r.Register(NewEnglishPack()) // Stub
    return r
}

func (r *Registry) Get(code string) (LanguagePack, error)
func (r *Registry) Register(pack LanguagePack)
```

### 9.3 French Implementation

```go
// internal/generator/languagepack/fr.go

type FrenchPack struct {
    tabooWords map[string]bool
    prompts    PromptTemplates
}

func (p *FrenchPack) Normalize(text string) string {
    // Remove accents: é→E, à→A, ç→C
    // Remove spaces and hyphens
    // Convert to uppercase A-Z
    return normalized
}

func (p *FrenchPack) IsTaboo(word string) bool {
    return p.tabooWords[strings.ToUpper(word)]
}
```

### 9.4 Adding a New Language

To add a new language (e.g., Spanish):

1. Create `internal/generator/languagepack/es.go`:

```go
type SpanishPack struct {
    tabooWords map[string]bool
    prompts    PromptTemplates
}

func NewSpanishPack() *SpanishPack {
    return &SpanishPack{
        tabooWords: loadSpanishTaboo(),
        prompts: PromptTemplates{
            ThemeGeneration: "Genera un tema de crucigrama...",
            SlotCandidates:  "Propón palabras de {length} letras...",
            ClueGeneration:  "Escribe definiciones para...",
        },
    }
}
```

2. Register in `DefaultRegistry()`:

```go
r.Register(NewSpanishPack())
```

3. Add normalization rules for Spanish characters (ñ, accents).

---

## 10. Flutter App

### 10.1 Architecture Overview

```
flutter_app/lib/
├── main.dart                 # Entry point with ProviderScope
├── app.dart                  # Material app, routing, themes
├── core/
│   ├── api/
│   │   ├── api_client.dart   # Dio HTTP wrapper
│   │   ├── api_endpoints.dart
│   │   └── api_exceptions.dart
│   ├── config/
│   │   └── app_config.dart   # Environment configuration
│   └── providers/
│       └── core_providers.dart
└── features/puzzle/
    ├── data/
    │   ├── models/           # Freezed immutable models
    │   │   ├── puzzle_model.dart
    │   │   ├── cell_model.dart
    │   │   └── clue_model.dart
    │   └── repositories/
    │       └── puzzle_repository.dart
    ├── providers/
    │   └── puzzle_providers.dart  # Riverpod state management
    └── presentation/
        ├── screens/
        │   └── home_screen.dart
        └── widgets/
            └── puzzle_card.dart
```

### 10.2 State Management (Riverpod)

```dart
// lib/features/puzzle/providers/puzzle_providers.dart

@riverpod
PuzzleRepository puzzleRepository(PuzzleRepositoryRef ref) {
  final client = ref.watch(apiClientProvider);
  return PuzzleRepository(client);
}

@riverpod
Future<Puzzle> dailyPuzzle(DailyPuzzleRef ref) async {
  final repo = ref.watch(puzzleRepositoryProvider);
  return repo.getDailyPuzzle(language: 'fr');
}
```

### 10.3 Data Models (Freezed)

```dart
// lib/features/puzzle/data/models/puzzle_model.dart

@freezed
class Puzzle with _$Puzzle {
  const factory Puzzle({
    required String id,
    required String date,
    required String language,
    required String title,
    required String author,
    required int difficulty,
    required String status,
    required List<List<Cell>> grid,
    required Clues clues,
  }) = _Puzzle;

  factory Puzzle.fromJson(Map<String, dynamic> json) => _$PuzzleFromJson(json);
}

// lib/features/puzzle/data/models/cell_model.dart

@freezed
class Cell with _$Cell {
  const Cell._();

  const factory Cell({
    required String type,     // "letter", "block", or "clue" (mots fléchés)
    String? solution,         // A-Z for letter cells
    int? number,              // Clue number (mots croisés style)
    String? clue,             // Definition text for clue cells
    String? arrow,            // "across" or "down" for clue cells
  }) = _Cell;

  bool get isLetter => type == 'letter';
  bool get isBlock => type == 'block';
  bool get isClue => type == 'clue';
  bool get isArrowDown => arrow == 'down';
}
```

### 10.4 Mots Fléchés Support

The app supports the **mots fléchés** (arrow words) format popular in French puzzles:

```
┌─────────────┬────┬────┬────┬────┬────┐
│  CAPITALE   │ P  │ A  │ R  │ I  │ S  │
│     →       │    │    │    │    │    │
├─────────────┼────┼────┼────┼────┼────┤
│  CHALEUR    │ C  │ H  │ A  │ U  │ D  │
│     →       │    │    │    │    │    │
└─────────────┴────┴────┴────┴────┴────┘
```

Clue cells contain:
- **Definition text**: The clue/hint for the word
- **Arrow indicator**: Points right (→) for across or down (↓) for down

This is rendered with a warm newspaper aesthetic:
- Paper background: `#FDF8F0`
- Tan clue cells: `#FCE4BC`
- Cyan selection: `#7FDBDB`
- Gray arrows: `#999999`

### 10.5 API Client

```dart
// lib/core/api/api_client.dart

class ApiClient {
  final Dio _dio;

  ApiClient(AppConfig config) : _dio = Dio(BaseOptions(
    baseUrl: config.apiBaseUrl,
    connectTimeout: Duration(seconds: 10),
    receiveTimeout: Duration(seconds: 10),
  ));

  Future<T> get<T>(String path, T Function(dynamic) fromJson) async {
    final response = await _dio.get(path);
    return fromJson(response.data);
  }
}
```

### 10.6 UI Components

The app uses Material 3 design with French localization:

```dart
// lib/features/puzzle/presentation/widgets/puzzle_card.dart

class PuzzleCard extends StatelessWidget {
  final Puzzle puzzle;
  final VoidCallback onPlay;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Column(
        children: [
          Text(puzzle.title),
          DifficultyIndicator(level: puzzle.difficulty),
          ElevatedButton(
            onPressed: onPlay,
            child: Text('Jouer'),
          ),
        ],
      ),
    );
  }
}
```

---

## 11. Deployment

### 11.1 Docker Compose Setup

```yaml
# compose.yaml

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile.api
    ports:
      - "8080:8080"
    volumes:
      - api-data:/app/data
    environment:
      - PORT=:8080
      - DATABASE_PATH=/app/data/puzzles.db
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3

  seed:
    image: curlimages/curl:latest
    depends_on:
      api:
        condition: service_healthy
    volumes:
      - ./seed:/seed:ro
    entrypoint: ["/bin/sh", "/seed/seed.sh"]

  web:
    build:
      context: ./flutter_app
    ports:
      - "3000:80"
    depends_on:
      seed:
        condition: service_completed_successfully

volumes:
  api-data:
```

### 11.2 Service Dependencies

```
┌─────────────────────────────────────────────────────────────────┐
│                    STARTUP ORDER                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. API Service starts                                          │
│     └─ Health check: GET /health                               │
│                                                                 │
│  2. Seed Service (after API healthy)                           │
│     └─ POSTs test puzzle to /admin/v1/puzzles                  │
│                                                                 │
│  3. Web Service (after Seed completes)                         │
│     └─ Flutter web app on port 3000                            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 11.3 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OPENAI_API_KEY` | (required) | OpenAI API key for generation |
| `PORT` | `:8080` | HTTP server bind address |
| `DATABASE_PATH` | `puzzles.db` | SQLite database file path |

### 11.4 API Dockerfile

```dockerfile
# Dockerfile.api

# Build stage
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o api ./cmd/api

# Runtime stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/schemas ./schemas
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["./api"]
```

### 11.5 Pre-commit Hooks

```yaml
# .pre-commit-config.yaml

repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json

  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
      - id: go-vet
      - id: go-mod-tidy
```

---

## 12. Key Files Reference

| Component | File Path |
|-----------|-----------|
| **Orchestrator** | `internal/generator/orchestrator.go` |
| **LLM Client** | `internal/generator/llm/client.go` |
| **OpenAI Impl** | `internal/generator/llm/openai.go` |
| **Theme Generator** | `internal/generator/theme/theme.go` |
| **Candidate Generator** | `internal/generator/theme/candidates.go` |
| **Grid Solver** | `internal/generator/fill/solver.go` |
| **Lexicon** | `internal/generator/fill/lexicon.go` |
| **Slot Discovery** | `internal/generator/fill/slot.go` |
| **Clue Generator** | `internal/generator/clue/clue.go` |
| **QA Scorer** | `internal/generator/qa/scorer.go` |
| **French Pack** | `internal/generator/languagepack/fr.go` |
| **Domain Types** | `internal/domain/types.go` |
| **API Routes** | `internal/api/routes.go` |
| **API Handlers** | `internal/api/handlers.go` |
| **Admin Handlers** | `internal/api/admin.go` |
| **Store Interface** | `internal/store/store.go` |
| **SQLite Store** | `internal/store/sqlite.go` |
| **API Main** | `cmd/api/main.go` |
| **Generator CLI** | `cmd/generate/main.go` |
| **Flutter App** | `flutter_app/lib/app.dart` |
| **Puzzle Model** | `flutter_app/lib/features/puzzle/data/models/puzzle_model.dart` |
| **Puzzle Providers** | `flutter_app/lib/features/puzzle/providers/puzzle_providers.dart` |
| **Docker Compose** | `compose.yaml` |
| **API Dockerfile** | `Dockerfile.api` |

---

## Appendix: Quick Start Commands

```bash
# Run API locally
go run ./cmd/api

# Generate a puzzle
OPENAI_API_KEY=sk-... go run ./cmd/generate -lang fr -difficulty 3

# Run tests
go test ./...

# Start full stack with Docker
docker compose up --build

# Seed the database
curl -X POST http://localhost:8080/admin/v1/puzzles \
  -H "Content-Type: application/json" \
  -d @seed/puzzle.json
```
