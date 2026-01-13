# Les Mots d'Atche - Executive Summary

## What It Is

Les Mots d'Atche is a French crossword puzzle application that uses AI to assist in creating high-quality puzzles daily. It combines the creativity of large language models (like GPT-4) with reliable algorithms to ensure every puzzle is solvable and engaging. Users can play puzzles on mobile (iOS/Android) or web with offline support.

---

## How It Works

The system uses a hybrid approach: **AI handles creativity, algorithms guarantee correctness.**

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│   AI picks      │ ──► │   Algorithm     │ ──► │   AI writes     │
│   theme & words │     │   fills grid    │     │   clues         │
│                 │     │                 │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
  "Cinema theme"          Guaranteed              "Definition for
   + word list           solvable grid             each word"
```

This hybrid design ensures puzzles are both creative and reliable.

---

## Key Capabilities

| Capability | Description |
|------------|-------------|
| **Daily Puzzles** | Fresh content generated automatically every day |
| **5 Difficulty Levels** | From casual (1) to expert (5) solvers |
| **Quality Assurance** | Automatic safety filters and quality scoring |
| **Offline Play** | Progress saved locally, no internet required |
| **French-First** | Optimized for French language, English-ready |
| **Cross-Platform** | iOS, Android, and web from a single codebase |

---

## Technical Highlights

| Component | Technology | Purpose |
|-----------|------------|---------|
| Backend | Go | Fast, reliable API server |
| AI Integration | OpenAI GPT-4 | Theme/clue generation |
| Database | SQLite | Lightweight, no server needed |
| Mobile App | Flutter | Cross-platform (iOS/Android/Web) |
| Deployment | Docker | Easy hosting and scaling |

---

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────────────┐
│                        USERS                                │
│         iOS App    Android App    Web Browser               │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      REST API                               │
│              (Puzzles, Health, Admin)                       │
└─────────────────────────┬───────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          │               │               │
          ▼               ▼               ▼
    ┌───────────┐   ┌───────────┐   ┌───────────┐
    │ Generator │   │  Database │   │    QA     │
    │ Pipeline  │   │  (SQLite) │   │  Scoring  │
    └─────┬─────┘   └───────────┘   └───────────┘
          │
          ▼
    ┌───────────┐
    │  OpenAI   │
    │   API     │
    └───────────┘
```

---

## Generation Pipeline

The puzzle generator follows a 7-step process:

| Step | Type | What It Does |
|------|------|--------------|
| 1. Theme | AI | Picks a theme with keywords |
| 2. Template | Algorithm | Creates grid layout |
| 3. Slots | Algorithm | Identifies word positions |
| 4. Candidates | AI | Suggests themed words |
| 5. Fill | Algorithm | Fills grid (guaranteed valid) |
| 6. Clues | AI | Writes puzzle clues |
| 7. QA | Algorithm | Scores and validates |

---

## Current Status

| Status | Feature |
|--------|---------|
| **Complete** | French puzzle generation |
| **Complete** | Mobile app (iOS/Android) |
| **Complete** | Daily puzzle delivery |
| **Complete** | Quality assurance system |
| **Complete** | Docker deployment |
| **Ready** | English language support (infrastructure in place) |
| **Planned** | Archive browsing |
| **Planned** | User accounts & progress sync |

---

## Key Metrics

- **Grid Size**: 13x13 (French standard)
- **Generation Time**: ~30-60 seconds per puzzle
- **Quality Threshold**: 60% minimum score to publish
- **Retry Attempts**: Up to 3 per generation
- **Languages**: French (active), English (ready)

---

## Getting Started

**Run locally:**
```bash
docker compose up --build
# API: http://localhost:8080
# Web: http://localhost:3000
```

**Generate a puzzle:**
```bash
OPENAI_API_KEY=sk-... go run ./cmd/generate -lang fr -difficulty 3
```

---

## Contact & Resources

- **Full Technical Docs**: See `docs/ARCHITECTURE.md`
- **API Endpoints**: `GET /v1/puzzles/daily?language=fr`
- **Repository**: Contains Go backend + Flutter app
