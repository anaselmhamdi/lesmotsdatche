#!/usr/bin/env python3
"""Command-line interface for crossword puzzle generation."""

import argparse
import json
import os
import sys
from datetime import datetime

from dotenv import load_dotenv

# Add the package to path
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from generator.orchestrator import GenerateRequest, Orchestrator
from language.french import FrenchPack
from llm.client import LLMClient, LLMConfig


def main():
    """Main entry point for CLI."""
    load_dotenv()

    parser = argparse.ArgumentParser(
        description="Generate French crossword puzzles",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument(
        "--date",
        default=datetime.now().strftime("%Y-%m-%d"),
        help="Puzzle date (YYYY-MM-DD)",
    )
    parser.add_argument(
        "--lang",
        default="fr",
        choices=["fr"],
        help="Language code",
    )
    parser.add_argument(
        "--difficulty",
        type=int,
        default=3,
        choices=[1, 2, 3, 4, 5],
        help="Difficulty level (1-5)",
    )
    parser.add_argument(
        "--max-size",
        type=int,
        default=10,
        help="Maximum grid dimension",
    )
    parser.add_argument(
        "--output",
        "-o",
        help="Output file path (default: stdout)",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Verbose output",
    )
    parser.add_argument(
        "--model",
        default="gpt-4o",
        help="OpenAI model to use",
    )

    args = parser.parse_args()

    # Check for API key
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        print("Error: OPENAI_API_KEY environment variable not set", file=sys.stderr)
        sys.exit(1)

    if args.verbose:
        print(f"Generating puzzle for {args.date}...", file=sys.stderr)
        print(f"  Language: {args.lang}", file=sys.stderr)
        print(f"  Difficulty: {args.difficulty}/5", file=sys.stderr)
        print(f"  Max size: {args.max_size}x{args.max_size}", file=sys.stderr)
        print(f"  Model: {args.model}", file=sys.stderr)

    # Initialize components
    llm_config = LLMConfig(model=args.model)
    llm_client = LLMClient(api_key=api_key, config=llm_config)
    orchestrator = Orchestrator(llm_client, FrenchPack())

    # Generate puzzle
    request = GenerateRequest(
        date=args.date,
        language=args.lang,
        difficulty=args.difficulty,
        max_size=args.max_size,
    )

    if args.verbose:
        print("Starting generation pipeline...", file=sys.stderr)

    result = orchestrator.generate(request)

    if not result.success:
        print(f"Error: {result.error}", file=sys.stderr)
        sys.exit(1)

    # Convert to JSON
    puzzle_dict = result.bundle.puzzle.model_dump(mode="json")
    report_dict = result.bundle.report.model_dump(mode="json")
    output = {
        "puzzle": puzzle_dict,
        "report": report_dict,
    }

    json_output = json.dumps(output, indent=2, ensure_ascii=False)

    # Output
    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(json_output)
        if args.verbose:
            print(f"Puzzle written to {args.output}", file=sys.stderr)
    else:
        print(json_output)

    if args.verbose:
        rows, cols = result.bundle.puzzle.grid_dimensions()
        print(f"\nGeneration complete!", file=sys.stderr)
        print(f"  Grid size: {rows}x{cols}", file=sys.stderr)
        print(f"  Clues: {len(result.bundle.puzzle.clues.across)} across, {len(result.bundle.puzzle.clues.down)} down", file=sys.stderr)
        print(f"  Fill score: {result.bundle.report.fill_score}/100", file=sys.stderr)
        print(f"  Freshness score: {result.bundle.report.freshness_score}/100", file=sys.stderr)


if __name__ == "__main__":
    main()
