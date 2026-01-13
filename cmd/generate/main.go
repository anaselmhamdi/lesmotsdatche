// Command generate creates crossword puzzles using the LLM-assisted generator.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"lesmotsdatche/internal/generator"
	"lesmotsdatche/internal/generator/fill"
	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
	"lesmotsdatche/internal/generator/theme"
)

func main() {
	// Load .env file if present (silently ignore if not found)
	_ = godotenv.Load()

	// Parse flags
	date := flag.String("date", time.Now().Format("2006-01-02"), "Target date (YYYY-MM-DD)")
	language := flag.String("lang", "fr", "Language code (fr, en)")
	difficulty := flag.Int("difficulty", 3, "Target difficulty (1-5)")
	maxSize := flag.Int("max-size", 12, "Max grid dimension (grid built around words)")
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
		fmt.Fprintf(os.Stderr, "Generating mots fléchés for %s in %s (max %dx%d, difficulty %d)\n",
			*date, langPack.Name(), *maxSize, *maxSize, *difficulty)
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

	// Create orchestrator with word-first approach
	config := generator.DefaultConfig()
	config.MaxAttempts = *maxAttempts
	config.Timeout = *timeout
	config.TargetDifficulty = *difficulty
	config.GridSize = [2]int{*maxSize, *maxSize} // Max bounds for word-first construction

	orch := generator.NewOrchestrator(validatingClient, langPack, baseLexicon, config)

	// Generate puzzle
	ctx := context.Background()
	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Starting word-first generation with model %s...\n", *model)
	}

	start := time.Now()
	result, err := orch.Generate(ctx, generator.GenerateRequest{
		Date:     *date,
		Language: *language,
		GridRows: *maxSize, // Max bounds, actual size determined by words
		GridCols: *maxSize,
		Constraints: theme.ThemeConstraints{
			Difficulty: *difficulty,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Generation failed: %v\n", err)
		// Print traces for debugging
		if *verbose {
			fmt.Fprintln(os.Stderr, "\nLLM Traces:")
			for i, trace := range validatingClient.Traces() {
				fmt.Fprintf(os.Stderr, "  [%d] Attempt %d:\n", i+1, trace.Attempt)
				fmt.Fprintf(os.Stderr, "      Prompt: %s\n", truncate(trace.Request.Prompt, 100))
				fmt.Fprintf(os.Stderr, "      Response: %s\n", truncate(trace.Response.Content, 200))
				if trace.Error != "" {
					fmt.Fprintf(os.Stderr, "      Error: %s\n", trace.Error)
				}
			}
		}
		os.Exit(1)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Generation completed in %v\n", time.Since(start))
		fmt.Fprintf(os.Stderr, "Theme: %s\n", result.Theme.Title)
		fmt.Fprintf(os.Stderr, "Grid size: %dx%d\n", len(result.Puzzle.Grid), len(result.Puzzle.Grid[0]))
		fmt.Fprintf(os.Stderr, "Words placed: %d\n", len(result.FillResult.Words))
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

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
