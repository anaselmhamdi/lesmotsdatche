// Command generate creates crossword puzzles using the LLM-assisted generator.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"lesmotsdatche/internal/generator"
	"lesmotsdatche/internal/generator/fill"
	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
	"lesmotsdatche/internal/generator/theme"
)

func main() {
	// Parse flags
	date := flag.String("date", time.Now().Format("2006-01-02"), "Target date (YYYY-MM-DD)")
	language := flag.String("lang", "fr", "Language code (fr, en)")
	difficulty := flag.Int("difficulty", 3, "Target difficulty (1-5)")
	output := flag.String("output", "", "Output file (default: stdout)")
	apiKey := flag.String("api-key", "", "OpenAI API key (or set OPENAI_API_KEY env)")
	model := flag.String("model", "gpt-4o", "LLM model to use")
	timeout := flag.Duration("timeout", 5*time.Minute, "Generation timeout")
	maxAttempts := flag.Int("max-attempts", 3, "Maximum generation attempts")
	verbose := flag.Bool("verbose", false, "Verbose output")

	flag.Parse()

	// Get API key
	key := *apiKey
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}
	if key == "" {
		fmt.Fprintln(os.Stderr, "Error: OpenAI API key required (use -api-key or set OPENAI_API_KEY)")
		os.Exit(1)
	}

	// Get language pack
	registry := languagepack.DefaultRegistry()
	langPack, ok := registry.Get(*language)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Unknown language: %s\n", *language)
		os.Exit(1)
	}
	if !langPack.IsConfigured() {
		fmt.Fprintf(os.Stderr, "Warning: Language pack '%s' is not fully configured\n", *language)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Generating puzzle for %s in %s (difficulty %d)\n", *date, langPack.Name(), *difficulty)
	}

	// Create LLM client
	openaiClient := llm.NewOpenAIClient(llm.OpenAIConfig{
		APIKey:  key,
		Model:   *model,
		Timeout: *timeout,
	})
	validatingClient := llm.NewValidatingClient(openaiClient, llm.DefaultConfig())

	// Create base lexicon
	baseLexicon := fill.SampleFrenchLexicon()

	// Create orchestrator
	config := generator.DefaultConfig()
	config.MaxAttempts = *maxAttempts
	config.Timeout = *timeout
	config.TargetDifficulty = *difficulty

	orch := generator.NewOrchestrator(validatingClient, langPack, baseLexicon, config)

	// Generate puzzle
	ctx := context.Background()
	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	if *verbose {
		fmt.Fprintln(os.Stderr, "Starting generation...")
	}

	start := time.Now()
	result, err := orch.Generate(ctx, generator.GenerateRequest{
		Date:     *date,
		Language: *language,
		Constraints: theme.ThemeConstraints{
			Difficulty: *difficulty,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Generation failed: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Generation completed in %v\n", time.Since(start))
		fmt.Fprintf(os.Stderr, "Theme: %s\n", result.Theme.Title)
		fmt.Fprintf(os.Stderr, "QA Score: %.2f\n", result.QAScore.Overall)
		fmt.Fprintf(os.Stderr, "Stats: %d attempts, %v fill time, %v clue time\n",
			result.Stats.Attempts, result.Stats.FillTime, result.Stats.ClueTime)
	}

	// Output result
	jsonData, err := json.MarshalIndent(result.Puzzle, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to encode puzzle: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		if err := os.WriteFile(*output, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to write output: %v\n", err)
			os.Exit(1)
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Puzzle written to %s\n", *output)
		}
	} else {
		fmt.Println(string(jsonData))
	}
}
