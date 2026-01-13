"""Tests for the lexicon module."""

import pytest

from generator.lexicon import MemoryLexicon, HybridLexicon


class TestMemoryLexicon:
    def test_add_words(self):
        lex = MemoryLexicon()
        lex.add_words(["HELLO", "WORLD", "TEST"])
        assert lex.contains("HELLO")
        assert lex.contains("WORLD")
        assert len(lex.words()) == 3

    def test_match_pattern(self):
        lex = MemoryLexicon(["HELLO", "HELPS", "WORLD"])
        matches = lex.match("HEL..")
        assert "HELLO" in matches
        assert "HELPS" in matches
        assert "WORLD" not in matches

    def test_words_by_length(self):
        lex = MemoryLexicon(["HI", "HELLO", "WORLD", "GO"])
        assert len(lex.words_by_length(2)) == 2
        assert len(lex.words_by_length(5)) == 2

    def test_normalize_on_add(self):
        lex = MemoryLexicon()
        lex.add_words(["CAFÉ", "résumé"])
        assert lex.contains("CAFE")
        assert lex.contains("RESUME")


class TestHybridLexicon:
    def test_primary_and_fallback(self):
        lex = HybridLexicon(["FALLBACK", "WORD"])
        lex.set_primary_words(["PRIMARY", "FIRST"])

        # Both should be accessible
        assert lex.contains("PRIMARY")
        assert lex.contains("FALLBACK")

    def test_match_prefers_primary(self):
        lex = HybridLexicon(["TEST"])
        lex.set_primary_words(["TEST", "BEST"])

        matches = lex.match("TEST")
        # Primary should come first
        assert matches[0] == "TEST"

    def test_primary_words_separate(self):
        lex = HybridLexicon(["FALLBACK"])
        lex.set_primary_words(["PRIMARY"])

        assert "PRIMARY" in lex.primary_words()
        assert "PRIMARY" not in lex.fallback_words()
        assert "FALLBACK" in lex.fallback_words()
