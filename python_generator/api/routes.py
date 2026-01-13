"""FastAPI routes for the crossword puzzle API."""

import os
from datetime import datetime
from typing import Optional

from fastapi import FastAPI, HTTPException, Query
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

from domain.types import DraftBundle, Puzzle, PuzzleStatus
from generator.orchestrator import GenerateRequest, Orchestrator
from language.french import FrenchPack
from llm.client import LLMClient, LLMConfig


# In-memory store for simplicity (replace with database in production)
_puzzle_store: dict[str, Puzzle] = {}


class GenerateRequestBody(BaseModel):
    """Request body for puzzle generation."""

    date: str
    language: str = "fr"
    difficulty: int = 3
    max_size: int = 10


class StatusUpdateBody(BaseModel):
    """Request body for status update."""

    status: PuzzleStatus


def create_app() -> FastAPI:
    """Create the FastAPI application."""
    app = FastAPI(
        title="Les Mots d'Atché API",
        description="French crossword puzzle generation API",
        version="1.0.0",
    )

    # CORS middleware for browser requests
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],  # In production, restrict to specific origins
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # Initialize LLM client
    api_key = os.getenv("OPENAI_API_KEY")
    llm_config = LLMConfig(
        model=os.getenv("OPENAI_MODEL", "gpt-4o"),
    )
    llm_client = LLMClient(api_key=api_key, config=llm_config)
    orchestrator = Orchestrator(llm_client, FrenchPack())

    @app.get("/health")
    async def health_check():
        """Health check endpoint."""
        return {"status": "ok", "timestamp": datetime.now().isoformat()}

    @app.get("/v1/puzzles/daily")
    async def get_daily_puzzle(language: str = Query(default="fr")):
        """Get today's puzzle."""
        today = datetime.now().strftime("%Y-%m-%d")

        # Look for a published puzzle for today
        for puzzle in _puzzle_store.values():
            if puzzle.date == today and puzzle.language == language and puzzle.status == PuzzleStatus.PUBLISHED:
                return puzzle

        raise HTTPException(status_code=404, detail="No puzzle available for today")

    @app.get("/v1/puzzles/{puzzle_id}")
    async def get_puzzle(puzzle_id: str):
        """Get puzzle by ID."""
        if puzzle_id not in _puzzle_store:
            raise HTTPException(status_code=404, detail="Puzzle not found")
        return _puzzle_store[puzzle_id]

    @app.post("/admin/v1/generate")
    async def generate_puzzle(request: GenerateRequestBody) -> DraftBundle:
        """Generate a new puzzle."""
        gen_request = GenerateRequest(
            date=request.date,
            language=request.language,
            difficulty=request.difficulty,
            max_size=request.max_size,
        )

        result = orchestrator.generate(gen_request)

        if not result.success:
            raise HTTPException(status_code=500, detail=result.error or "Generation failed")

        return result.bundle

    @app.post("/admin/v1/puzzles")
    async def store_puzzle(puzzle: Puzzle):
        """Store a puzzle."""
        _puzzle_store[puzzle.id] = puzzle
        return {"id": puzzle.id, "status": "stored"}

    @app.patch("/admin/v1/puzzles/{puzzle_id}/status")
    async def update_puzzle_status(puzzle_id: str, body: StatusUpdateBody):
        """Update puzzle status."""
        if puzzle_id not in _puzzle_store:
            raise HTTPException(status_code=404, detail="Puzzle not found")

        puzzle = _puzzle_store[puzzle_id]
        puzzle.status = body.status

        if body.status == PuzzleStatus.PUBLISHED:
            puzzle.published_at = datetime.now()

        return {"id": puzzle_id, "status": puzzle.status}

    return app


# For running with uvicorn
app = create_app()
