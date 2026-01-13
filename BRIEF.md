# Modern Crossword (FR-first) — LLM-Generation-First Spec (Go + Flutter)

> **Goal:** Ship a modern-feeling crossword app starting in **French**, with **LLM-assisted generation** as the core differentiator. Architecture must make it easy to add **English** later with minimal refactors (language packs, lexicons, clue styles, normalization rules).

---

## 1) Product Overview

### 1.1 Core promise
- A fast, offline-friendly crossword player.
- Daily + archive puzzles.
- High-quality modern references in clueing and fill.
- A generation pipeline that uses LLMs for **theme/entries + clue writing + difficulty calibration**, while deterministic code guarantees **grid correctness**.

### 1.2 MVP user stories
1. Open app → play **Daily** immediately.
2. Archive browsing with filters: date, difficulty, tags (theme/ref).
3. Crossword player: tap/type, smooth navigation, clue focus bar.
4. Offline progress save/resume.
5. Optional: check/reveal letter/word/grid with limits.
6. Admin workflow: generate draft → review → publish.

### 1.3 Non-goals (MVP)
- No public puzzle submission.
- No automatic publishing of LLM-generated puzzles without review.
- No heavy auth; admin can be basic auth.

---

## 2) Internationalization Strategy (FR-first, EN-ready)

### 2.1 Language as a first-class concept
- Every puzzle has `language` (`fr` for MVP, `en` future).
- Every generator run is parameterized by language:
  - lexicon source
  - normalization rules
  - taboo lists
  - clue style guidelines
  - modern reference seed set
- Client UI strings use Flutter localization (`intl`), with FR as default; EN can be added via ARB files.

### 2.2 Character / normalization policy (MVP)
- **Grid uses A–Z only** (no diacritics in cells).
- Answers are normalized:
  - FR: strip accents, remove spaces/hyphens/apostrophes for grid, keep display in clues.
  - EN: similar (no diacritics usually).
- Future enhancement: optional accented cell support by language pack.

---

## 3) System Architecture

### 3.1 Components
1. **API (Go)**  
   Serves published puzzles + metadata, and admin endpoints for drafts/publishing.
2. **Generator (Go)** *(primary emphasis)*  
   Hybrid generator: deterministic fill solver + LLM-assisted theme/entries and clue writing.
3. **DB**  
   Postgres (prod) + SQLite (dev ok). Store puzzles, drafts, scores, review notes.
4. **Flutter app**  
   Offline-first player, modern UX, easily extensible to English.

### 3.2 High-level flow
- Admin triggers generation (or scheduled job later).
- Generator:
  - gets theme & candidates via LLM
  - fills grid with deterministic solver (with LLM used to propose candidates for hard slots)
  - generates clue variants via LLM
  - runs QA scoring + safety checks
  - stores as `draft`
- Admin reviews and publishes to `published`
- Client fetches daily/archived puzzles and stores locally; progress saved locally.

---

## 4) Data Model

### 4.1 Core entities

#### Puzzle
- `id` UUID
- `date` YYYY-MM-DD (daily key)
- `language` `"fr"` | `"en"`
- `title` string
- `author` string
- `difficulty` 1..5
- `status` `"draft"` | `"published"` | `"archived"`
- `grid` 2D array of `Cell`
- `clues` { `across`: Clue[], `down`: Clue[] }
- `metadata`:
  - `theme_tags`: string[]
  - `reference_tags`: string[]
  - `notes`: string
  - `freshness_score`: 0..100
- timestamps: `created_at`, `published_at`

#### Cell
- `type`: `"letter"` | `"block"` | `"clue"` (mots fléchés)
- `solution`: `"A"`..`"Z"` (for letter cells)
- `number`: int? (computed, for mots croisés style)
- `clue_across`: string? (definition for across direction →)
- `clue_down`: string? (definition for down direction ↓)
- (Client-side only) `entry`: `"A"`..`"Z"`?

**Mots Fléchés Format:**
In mots fléchés (arrow words), clue cells are embedded directly in the grid. Each clue cell contains one or two definitions with arrows indicating the answer direction.

**Single clue cells:** Have either `clue_across` OR `clue_down` set.
**Split clue cells:** Have BOTH `clue_across` AND `clue_down` set, with a horizontal divider separating the two definitions. This allows one cell to define both an across and down word, making grid construction more flexible.

#### Clue
- `id` UUID
- `direction`: `"across"` | `"down"`
- `number`: int
- `prompt`: string (localized to puzzle language)
- `answer`: string (normalized A–Z for grid)
- `start`: `{row:int, col:int}`
- `length`: int
- `reference_tags`: string[]
- `reference_year_range`: `[int,int]`
- `difficulty`: 1..5
- `ambiguity_notes`: string?

#### DraftReport (stored with draft)
- `fill_score`: 0..100
- `clue_score`: 0..100
- `freshness_score`: 0..100
- `risk_flags`: string[] (e.g., `AMBIGUOUS_CLUE`, `TOO_MANY_PROPER_NOUNS`, `OBSCURE_ENTRY`)
- `slot_failures`: { pattern:string, length:int, attempts:int }[]
- `language_checks`: { taboo_hits:int, proper_nouns:int, avg_word_freq:float }
- `llm_trace_ref`: string? (optional; store prompt/response separately, redact secrets)

---

## 5) Puzzle File Format & Schemas

### 5.1 Canonical JSON
- Define JSONSchema for:
  - `Puzzle` (published)
  - `DraftPuzzleBundle` = Puzzle + DraftReport + optional alternate clue sets
  - `LLMThemeResponse`, `LLMSlotCandidatesResponse`, `LLMClueResponse`

### 5.2 Deterministic numbering
- Number assigned where a cell starts an across or down entry:
  - across start if left is block/outside and right is letter
  - down start if up is block/outside and down is letter
- Numbering stable row-major.

### 5.3 Validation rules (must)
- Rectangular grid, min size (e.g., 7x7 for MVP, allow 15x15 later).
- All clue answers match grid letters and lengths.
- Crossings consistent.
- No invalid characters in normalized answers.
- Ensure each non-block cell belongs to at least one entry.
- Optional: grid connectivity check (one island).

---

## 6) LLM-First Generator Requirements (Core Differentiator)

### 6.1 Responsibilities split
**LLM does:**
- Theme + candidate entries list
- Candidate suggestions for hard slots (pattern-based)
- Clue writing (multiple variants)
- Optional judging of clue fairness/ambiguity

**Deterministic code does:**
- Grid creation/selection
- Fill solving and constraint satisfaction
- Validation/canonicalization
- QA heuristics + safety filters
- Publishing gates

### 6.2 Language packs (FR now, EN later)
A language pack includes:
- normalization rules
- lexicon path and frequency metadata
- taboo list and “avoid list”
- prompt templates for LLM (theme + clue style)
- reference seed dataset

### 6.3 LLM output contracts (strict JSON)
All LLM outputs MUST:
- validate against JSONSchema
- be retried via “repair prompt” on failure
- be stored for debugging (redacted secrets)

#### Theme & candidate entries response
```json
{
  "theme_title": "string",
  "theme_description": "string",
  "theme_tags": ["string"],
  "candidates": [
    {
      "answer": "string",
      "reference_tags": ["string"],
      "reference_year_range": [2018, 2026],
      "difficulty": 1,
      "notes": "string"
    }
  ]
}

#### Slot-candidates response

-   Input: `{language, pattern, length, desired_tags[], difficulty_target}`

-   Output: `{candidates:[{answer, tags, year_range, difficulty}]}`

#### Clue response

-   Input: `{language, answer, tags, difficulty_target, style_guide}`

-   Output: `{variants:[{prompt, difficulty, ambiguity_notes?}]}`

### 6.4 Fill solver (deterministic)

-   Constraint-based fill:

    -   choose next slot by most constrained (fewest candidates)

    -   score candidates by:

        -   frequency score

        -   theme fit (tags overlap)

        -   diversity penalty (avoid repeating tag types)

        -   crossing friendliness

-   When candidate set too small, call slot-candidates LLM helper.

### 6.5 QA & scoring

-   Flag ambiguous clues (heuristics: too short prompt, missing disambiguators, common homonyms).

-   Penalize:

    -   too many proper nouns vs difficulty

    -   rare words in easy puzzles

    -   repeated clue patterns

-   Freshness:

    -   boost entries with year_range overlapping 2018--current year

    -   ensure balanced modern/classic to avoid trivia dump

### 6.6 Safety

-   Avoid: private individuals, doxxing, harassment, slurs.

-   Maintain taboo lists per language pack.

-   If LLM suggests risky content, replace candidate and record flag.

### 6.7 Publishing rule

-   Drafts NEVER auto-publish in MVP.

-   Publishing requires admin action and passing validation.

* * * * *

7) Go API Spec
--------------

### 7.1 Public endpoints

-   `GET /v1/puzzles/daily?language=fr` → puzzle payload

-   `GET /v1/puzzles?language=fr&from=YYYY-MM-DD&to=...&tag=...&difficulty=...` → metadata list

-   `GET /v1/puzzles/{id}` → payload

Caching:

-   Puzzle payloads include ETag + Cache-Control.

-   Gzip enabled.

### 7.2 Admin endpoints (basic auth)

-   `POST /v1/admin/generate` → create a draft (sync MVP)

    -   body: `{language:"fr", size:15, difficulty_target:2, theme_tags?:[]}`

-   `GET /v1/admin/drafts?language=fr`

-   `GET /v1/admin/drafts/{id}`

-   `POST /v1/admin/drafts/{id}/publish`

-   `POST /v1/admin/drafts/{id}/reject` body `{reason:string, notes?:string}`

-   `POST /v1/admin/puzzles/validate` (validate uploaded JSON)

-   `POST /v1/admin/puzzles` (upload draft JSON)

* * * * *

8) Flutter App Spec (FR-first)
------------------------------

### 8.1 Tech choices

-   `go_router` for navigation

-   `flutter_riverpod` for state management

-   `drift` (SQLite) for offline persistence

-   `dio` for networking

-   localization: `intl` + ARB files

### 8.2 Screens

1.  Home:

    -   Daily puzzle card

    -   Continue last puzzle

    -   language toggle hidden for MVP (FR only), but architecture supports it

2.  Archive:

    -   list view

    -   filter chips: difficulty, theme tags

    -   search by title/date

3.  Player:

    -   grid

    -   clue list panel

    -   sticky current clue bar

    -   controls: direction toggle, reveal/check, timer, mistake mode

4.  Settings:

    -   dark mode toggle

    -   future: language selection (FR/EN)

### 8.3 Player UX requirements

-   Tap cell selects.

-   Typing fills and auto-advances.

-   Backspace clears then moves back.

-   Tap clue selects word and highlights.

-   Toggle direction by button or tapping active cell again.

-   Highlight:

    -   active cell

    -   active word

    -   crossing word

-   Offline autosave progress on every change.

### 8.4 Offline persistence schema (client)

-   `puzzles` table: id, language, date, title, payload_json, etag

-   `progress` table: puzzle_id, entries_json, elapsed_seconds, last_cell, reveal_counts, updated_at

* * * * *

9) Repo Structure
-----------------

`/cmd/api
/cmd/generator
/internal/
  api/            handlers + middleware
  domain/         puzzle structs, numbering, extraction
  validate/       jsonschema validation, canonicalization
  generator/
    languagepack/ fr.go, en.go (stub)
    llm/          client, retries, schema validation
    fill/         solver
    clue/         clue orchestration
    scoring/      QA + scoring
  store/          repos (postgres/sqlite)
/schemas/
/seed/
  fr_seed_modern.json
  en_seed_modern.json (stub)
/flutter_app/
  lib/
    app/
    features/
      home/
      archive/
      player/
      settings/
    data/
      api/
      db/
    l10n/
/docker-compose.yml`

* * * * *

10) Testing Requirements
------------------------

### 10.1 Go tests

-   Numbering correctness (golden).

-   Across/down extraction (golden).

-   Schema validation errors (table-driven).

-   Fill solver determinism (fixed seed).

-   LLM JSON parsing + repair retry logic (mock LLM).

### 10.2 Flutter tests

-   Widget tests for player navigation.

-   Persistence tests (save/resume).

-   Basic golden/snapshot tests for grid rendering (optional MVP).

* * * * *

11) Acceptance Criteria (MVP)
-----------------------------

-   Generate a draft crossword in FR via admin endpoint.

-   Draft includes QA report and multiple clue variants.

-   Publish flow promotes draft → public.

-   Flutter app can:

    -   fetch daily puzzle (FR),

    -   play smoothly,

    -   save progress offline,

    -   resume reliably.

-   Codebase supports adding EN by adding an English language pack + lexicon + UI strings, without changing core domain logic.

* * * * *