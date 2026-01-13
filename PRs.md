
---

## PR plan: small, testable tasks for the agent

Below is a breakdown into **PR-sized chunks** that should each be reviewable and testable. (Each PR should include tests and a short `RUN.md` update.)

### PR 1 — Go domain model + canonical puzzle format
- Add `Puzzle/Cell/Clue` structs
- Implement normalization (FR rules) + stubs for EN rules
- Implement deterministic numbering + across/down extraction
- Add golden fixtures + tests (`testdata/*.json`)

**Testable:** `go test ./...` passes; golden files stable.

---

### PR 2 — JSONSchema + validator
- Add JSONSchema files for `Puzzle` + `DraftBundle`
- Implement schema validation in Go
- Add clear validation error messages (path + reason)
- Add table-driven tests for invalid puzzles

**Testable:** validator rejects known-bad fixtures.

---

### PR 3 — Storage layer + migrations
- Add DB schema/migrations for puzzles, drafts, reports
- Implement repository interfaces (store/get/list)
- Add tests using SQLite in-memory (or testcontainers if you want)

**Testable:** CRUD tests pass.

---

### PR 4 — Public API endpoints (read-only)
- Implement `GET /v1/puzzles/daily`, `GET /v1/puzzles`, `GET /v1/puzzles/{id}`
- Add ETag + gzip + basic logging
- Add HTTP tests

**Testable:** curl endpoints locally; tests cover status codes + JSON.

---

### PR 5 — LLM client wrapper + strict JSON output handling
- Implement `llm.Client` abstraction (provider-agnostic)
- Implement:
  - JSONSchema validation for LLM outputs
  - retry + “repair prompt” strategy
  - redaction of secrets in traces
- Add mock LLM + tests for failure modes

**Testable:** unit tests simulate malformed JSON then repair success.

---

### PR 6 — Language pack system (FR implemented, EN stub)
- Implement `LanguagePack` interface
- Add `fr` pack: normalization, taboo list, prompt templates, seed path
- Add `en` stub pack: placeholders + tests verifying pack selection works

**Testable:** generator can run with `fr`; `en` compiles but returns “not configured” for lexicon.

---

### PR 7 — Deterministic fill solver (no LLM yet)
- Implement slot discovery + constraints
- Implement backtracking fill with scoring hooks
- Use a small lexicon fixture for tests
- Add determinism tests with fixed random seed

**Testable:** given template + lexicon fixture, solver fills grid consistently.

---

### PR 8 — LLM-assisted theme + slot candidate integration
- Implement theme generation call
- Merge LLM candidates into solver candidate sets
- Implement slot-helper call when solver is stuck
- Add tests using mock LLM returning fixed JSON

**Testable:** solver succeeds on a grid that fails without LLM candidates.

---

### PR 9 — LLM clue generation + multi-variant clues
- Implement clue generation for each answer
- Store multiple variants; select default per difficulty
- Add ambiguity flags
- Tests: mock LLM outputs → stored correctly; schema enforced

**Testable:** draft bundle contains variants and passes validator.

---

### PR 10 — QA scoring + safety filters
- Implement `DraftReport` scoring
- Implement taboo/proper noun heuristics (simple rules MVP)
- Add risk flags + tests

**Testable:** known “bad” entries flagged; report fields populated.

---

### PR 11 — Generator command + admin endpoints
- Add `cmd/generator` runnable (or generator package invoked by API)
- Implement admin endpoints:
  - generate draft
  - list drafts
  - publish/reject
- Add HTTP tests with mock LLM

**Testable:** end-to-end: generate → draft stored → publish → appears in public list.

---

### PR 12 — Flutter app scaffold + networking layer
- Create Flutter project under `/flutter_app`
- Add routing, Riverpod setup, dio client
- Implement API models + fetch daily puzzle
- Add basic home screen (FR strings)

**Testable:** `flutter test` passes; app runs and shows daily puzzle title.

---

### PR 13 — Flutter player MVP (grid + input + navigation)
- Grid rendering (letters/blocks/clue cells) + selection
- **Mots fléchés support:** clue cells with embedded definitions and arrow indicators
- Typing/backspace navigation rules
- Current clue bar (from selected cell's clue)
- Unit/widget tests for navigation logic
- French newspaper aesthetic theme (warm paper, tan clue cells, cyan selection)

**Testable:** widget tests for direction toggle and cell movement; visual rendering of mots fléchés grid.

---

### PR 14 — Offline persistence (Drift) + resume
- Add drift DB schema
- Cache puzzle payloads by ETag
- Save progress on each change
- Resume last played puzzle automatically
- Tests for persistence roundtrip

**Testable:** kill/restart simulation restores grid.

---

### PR 15 — Polish + extensibility hooks for English
- Add localization framework with ARB files (`fr` now, `en` placeholders)
- Add settings screen with future language selector (disabled)
- Dark mode toggle
- Final docs: `RUN.md`, `ARCHITECTURE.md`, `GENERATOR.md`

**Testable:** switching locale to EN (dev) changes UI strings without breaking.

---

If you want, I can also add a **“Definition of Done” template** for PRs (linting, formatting, minimum tests, example curl commands) tailored to Go + Flutter, but the above is already structured so a coding agent can execute it sequentially.
