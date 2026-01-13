// Package validate provides JSON schema and semantic validation for puzzles.
package validate

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"lesmotsdatche/internal/domain"
)

//go:embed schemas/*.json
var schemasFS embed.FS

var (
	puzzleSchema     *jsonschema.Schema
	draftBundleSchema *jsonschema.Schema
)

func init() {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft2020

	// Load puzzle schema
	puzzleData, err := schemasFS.ReadFile("schemas/puzzle.schema.json")
	if err != nil {
		panic(fmt.Sprintf("failed to read puzzle schema: %v", err))
	}
	if err := compiler.AddResource("puzzle.schema.json", strings.NewReader(string(puzzleData))); err != nil {
		panic(fmt.Sprintf("failed to add puzzle schema: %v", err))
	}

	puzzleSchema, err = compiler.Compile("puzzle.schema.json")
	if err != nil {
		panic(fmt.Sprintf("failed to compile puzzle schema: %v", err))
	}

	// Load draft bundle schema
	draftData, err := schemasFS.ReadFile("schemas/draft_bundle.schema.json")
	if err != nil {
		panic(fmt.Sprintf("failed to read draft bundle schema: %v", err))
	}
	if err := compiler.AddResource("draft_bundle.schema.json", strings.NewReader(string(draftData))); err != nil {
		panic(fmt.Sprintf("failed to add draft bundle schema: %v", err))
	}

	draftBundleSchema, err = compiler.Compile("draft_bundle.schema.json")
	if err != nil {
		panic(fmt.Sprintf("failed to compile draft bundle schema: %v", err))
	}
}

// ValidationError represents a single validation error with path context.
type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no errors"
	}
	var msgs []string
	for _, e := range ve {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "; ")
}

// ValidatePuzzleJSON validates puzzle JSON against the schema.
func ValidatePuzzleJSON(data []byte) ValidationErrors {
	var doc interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return ValidationErrors{{Path: "", Message: fmt.Sprintf("invalid JSON: %v", err)}}
	}

	if err := puzzleSchema.Validate(doc); err != nil {
		return schemaErrorToValidationErrors(err)
	}

	return nil
}

// ValidateDraftBundleJSON validates draft bundle JSON against the schema.
func ValidateDraftBundleJSON(data []byte) ValidationErrors {
	var doc interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return ValidationErrors{{Path: "", Message: fmt.Sprintf("invalid JSON: %v", err)}}
	}

	if err := draftBundleSchema.Validate(doc); err != nil {
		return schemaErrorToValidationErrors(err)
	}

	return nil
}

// schemaErrorToValidationErrors converts jsonschema errors to ValidationErrors.
func schemaErrorToValidationErrors(err error) ValidationErrors {
	var errors ValidationErrors

	switch e := err.(type) {
	case *jsonschema.ValidationError:
		errors = append(errors, extractValidationErrors(e)...)
	default:
		errors = append(errors, ValidationError{
			Path:    "",
			Message: err.Error(),
		})
	}

	return errors
}

func extractValidationErrors(ve *jsonschema.ValidationError) ValidationErrors {
	var errors ValidationErrors

	// Add the main error if it has a message
	if ve.Message != "" {
		errors = append(errors, ValidationError{
			Path:    ve.InstanceLocation,
			Message: ve.Message,
		})
	}

	// Recursively add child errors
	for _, cause := range ve.Causes {
		errors = append(errors, extractValidationErrors(cause)...)
	}

	return errors
}

// ValidatePuzzleSemantic performs semantic validation on a parsed puzzle.
// This catches errors that JSON Schema cannot express.
func ValidatePuzzleSemantic(p *domain.Puzzle) ValidationErrors {
	var errors ValidationErrors

	// Check grid is rectangular
	if len(p.Grid) > 0 {
		expectedCols := len(p.Grid[0])
		for i, row := range p.Grid {
			if len(row) != expectedCols {
				errors = append(errors, ValidationError{
					Path:    fmt.Sprintf("/grid/%d", i),
					Message: fmt.Sprintf("row has %d columns, expected %d", len(row), expectedCols),
				})
			}
		}
	}

	// Check grid size constraints (French standard: 10-16)
	const (
		MinGridSize = 10
		MaxGridSize = 16
	)
	rows, cols := p.GridDimensions()
	if rows < MinGridSize || rows > MaxGridSize || cols < MinGridSize || cols > MaxGridSize {
		errors = append(errors, ValidationError{
			Path:    "/grid",
			Message: fmt.Sprintf("grid must be %dx%d to %dx%d, got %dx%d", MinGridSize, MinGridSize, MaxGridSize, MaxGridSize, rows, cols),
		})
	}

	// Check all letter cells have valid solutions
	for r, row := range p.Grid {
		for c, cell := range row {
			if cell.IsLetter() {
				if len(cell.Solution) != 1 || cell.Solution[0] < 'A' || cell.Solution[0] > 'Z' {
					errors = append(errors, ValidationError{
						Path:    fmt.Sprintf("/grid/%d/%d/solution", r, c),
						Message: fmt.Sprintf("letter cell must have A-Z solution, got %q", cell.Solution),
					})
				}
			}
		}
	}

	// Validate clue answers match grid
	for i, clue := range p.Clues.Across {
		gridAnswer := extractAnswer(p.Grid, clue.Start, clue.Length, domain.DirectionAcross)
		if gridAnswer != clue.Answer {
			errors = append(errors, ValidationError{
				Path:    fmt.Sprintf("/clues/across/%d/answer", i),
				Message: fmt.Sprintf("answer %q doesn't match grid %q", clue.Answer, gridAnswer),
			})
		}
	}

	for i, clue := range p.Clues.Down {
		gridAnswer := extractAnswer(p.Grid, clue.Start, clue.Length, domain.DirectionDown)
		if gridAnswer != clue.Answer {
			errors = append(errors, ValidationError{
				Path:    fmt.Sprintf("/clues/down/%d/answer", i),
				Message: fmt.Sprintf("answer %q doesn't match grid %q", clue.Answer, gridAnswer),
			})
		}
	}

	// Check clue lengths match answer lengths
	for i, clue := range p.Clues.Across {
		if clue.Length != len(clue.Answer) {
			errors = append(errors, ValidationError{
				Path:    fmt.Sprintf("/clues/across/%d/length", i),
				Message: fmt.Sprintf("length %d doesn't match answer length %d", clue.Length, len(clue.Answer)),
			})
		}
	}

	for i, clue := range p.Clues.Down {
		if clue.Length != len(clue.Answer) {
			errors = append(errors, ValidationError{
				Path:    fmt.Sprintf("/clues/down/%d/length", i),
				Message: fmt.Sprintf("length %d doesn't match answer length %d", clue.Length, len(clue.Answer)),
			})
		}
	}

	// Check every non-block cell belongs to at least one entry
	cellCoverage := make(map[string]bool)
	for _, clue := range p.Clues.Across {
		for i := 0; i < clue.Length; i++ {
			key := fmt.Sprintf("%d,%d", clue.Start.Row, clue.Start.Col+i)
			cellCoverage[key] = true
		}
	}
	for _, clue := range p.Clues.Down {
		for i := 0; i < clue.Length; i++ {
			key := fmt.Sprintf("%d,%d", clue.Start.Row+i, clue.Start.Col)
			cellCoverage[key] = true
		}
	}

	for r, row := range p.Grid {
		for c, cell := range row {
			if cell.IsLetter() {
				key := fmt.Sprintf("%d,%d", r, c)
				if !cellCoverage[key] {
					errors = append(errors, ValidationError{
						Path:    fmt.Sprintf("/grid/%d/%d", r, c),
						Message: "letter cell is not part of any clue entry",
					})
				}
			}
		}
	}

	return errors
}

func extractAnswer(grid [][]Cell, start domain.Position, length int, dir domain.Direction) string {
	var answer strings.Builder
	for i := 0; i < length; i++ {
		var r, c int
		if dir == domain.DirectionAcross {
			r, c = start.Row, start.Col+i
		} else {
			r, c = start.Row+i, start.Col
		}
		if r >= 0 && r < len(grid) && c >= 0 && c < len(grid[r]) {
			answer.WriteString(grid[r][c].Solution)
		}
	}
	return answer.String()
}

// Cell is a local type alias for embedding compatibility
type Cell = domain.Cell

// ValidatePuzzle performs both schema and semantic validation.
func ValidatePuzzle(data []byte) ValidationErrors {
	// First validate schema
	schemaErrors := ValidatePuzzleJSON(data)
	if len(schemaErrors) > 0 {
		return schemaErrors
	}

	// Then parse and validate semantically
	var puzzle domain.Puzzle
	if err := json.Unmarshal(data, &puzzle); err != nil {
		return ValidationErrors{{Path: "", Message: fmt.Sprintf("failed to parse puzzle: %v", err)}}
	}

	return ValidatePuzzleSemantic(&puzzle)
}
