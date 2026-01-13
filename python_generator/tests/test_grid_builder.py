"""Tests for the grid builder module."""

import pytest

from domain.types import CellType
from generator.grid_builder import GridBuilder, BuilderConfig


class TestGridBuilder:
    def test_build_basic(self):
        config = BuilderConfig(max_rows=10, max_cols=10, seed=42)
        builder = GridBuilder(config)

        candidates = [
            "HELLO", "WORLD", "TESTS", "WORDS", "CROSS",
            "TABLE", "CHAIR", "HOUSE", "WATER", "MUSIC",
        ]

        result = builder.build(candidates)

        # Should have placed some words
        assert len(result.words) > 0
        assert len(result.grid) > 0

    def test_build_with_short_words(self):
        config = BuilderConfig(max_rows=10, max_cols=10, seed=42)
        builder = GridBuilder(config)

        # Include short words for gap filling
        candidates = [
            "CINEMA", "ACTEUR", "SCENE", "FILM", "ROLE",
            "DE", "LA", "LE", "UN", "EN",
        ]

        result = builder.build(candidates)
        assert result.success or len(result.words) >= 2

    def test_grid_has_letters_and_blocks(self):
        config = BuilderConfig(max_rows=10, max_cols=10, seed=42)
        builder = GridBuilder(config)

        candidates = ["HELLO", "WORLD", "TESTS", "CROSS", "WORDS"]
        result = builder.build(candidates)

        if result.success:
            has_letters = False
            has_blocks = False

            for row in result.grid:
                for cell in row:
                    if cell.type == CellType.LETTER:
                        has_letters = True
                    if cell.type == CellType.BLOCK:
                        has_blocks = True

            assert has_letters
            assert has_blocks

    def test_deterministic_with_seed(self):
        candidates = ["HELLO", "WORLD", "TESTS", "CROSS", "WORDS"]

        # Same seed should produce same result
        config1 = BuilderConfig(max_rows=10, max_cols=10, seed=12345)
        result1 = GridBuilder(config1).build(candidates)

        config2 = BuilderConfig(max_rows=10, max_cols=10, seed=12345)
        result2 = GridBuilder(config2).build(candidates)

        assert result1.words == result2.words
