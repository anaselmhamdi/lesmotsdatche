package domain

import "strings"

// ExtractSlots extracts skeleton clues (slots) from a numbered grid.
// It returns clues with positions and answers but no prompts.
// The grid must have been processed by AssignNumbers first.
func ExtractSlots(grid [][]Cell) Clues {
	if len(grid) == 0 {
		return Clues{}
	}

	rows := len(grid)
	cols := len(grid[0])

	var across []Clue
	var down []Clue

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cell := grid[row][col]

			// Only process numbered cells
			if cell.Number == 0 {
				continue
			}

			// Check for across entry
			if StartsAcross(grid, row, col) {
				clue := extractAcrossClue(grid, row, col, cols, cell.Number)
				across = append(across, clue)
			}

			// Check for down entry
			if StartsDown(grid, row, col) {
				clue := extractDownClue(grid, row, col, rows, cell.Number)
				down = append(down, clue)
			}
		}
	}

	return Clues{
		Across: across,
		Down:   down,
	}
}

// extractAcrossClue extracts an across clue starting at (row, col).
func extractAcrossClue(grid [][]Cell, row, col, cols, number int) Clue {
	var answer strings.Builder
	startCol := col

	for c := col; c < cols && grid[row][c].IsLetter(); c++ {
		answer.WriteString(grid[row][c].Solution)
	}

	return Clue{
		Direction: DirectionAcross,
		Number:    number,
		Answer:    answer.String(),
		Start:     Position{Row: row, Col: startCol},
		Length:    answer.Len(),
	}
}

// extractDownClue extracts a down clue starting at (row, col).
func extractDownClue(grid [][]Cell, row, col, rows, number int) Clue {
	var answer strings.Builder
	startRow := row

	for r := row; r < rows && grid[r][col].IsLetter(); r++ {
		answer.WriteString(grid[r][col].Solution)
	}

	return Clue{
		Direction: DirectionDown,
		Number:    number,
		Answer:    answer.String(),
		Start:     Position{Row: startRow, Col: col},
		Length:    answer.Len(),
	}
}

// GetCellsForClue returns the positions of all cells belonging to a clue.
func GetCellsForClue(clue Clue) []Position {
	cells := make([]Position, clue.Length)

	for i := 0; i < clue.Length; i++ {
		if clue.Direction == DirectionAcross {
			cells[i] = Position{Row: clue.Start.Row, Col: clue.Start.Col + i}
		} else {
			cells[i] = Position{Row: clue.Start.Row + i, Col: clue.Start.Col}
		}
	}

	return cells
}

// FindClueAt finds the clue(s) that contain the given cell position.
// Returns both across and down clues if the cell is at a crossing.
func FindCluesAt(clues Clues, pos Position) (across *Clue, down *Clue) {
	// Check across clues
	for i := range clues.Across {
		c := &clues.Across[i]
		if pos.Row == c.Start.Row &&
			pos.Col >= c.Start.Col &&
			pos.Col < c.Start.Col+c.Length {
			across = c
			break
		}
	}

	// Check down clues
	for i := range clues.Down {
		c := &clues.Down[i]
		if pos.Col == c.Start.Col &&
			pos.Row >= c.Start.Row &&
			pos.Row < c.Start.Row+c.Length {
			down = c
			break
		}
	}

	return
}

// ValidateCluesAgainstGrid checks that all clue answers match the grid.
// Returns a list of mismatches (empty if all valid).
func ValidateCluesAgainstGrid(grid [][]Cell, clues Clues) []string {
	var errors []string

	for _, clue := range clues.Across {
		gridAnswer := extractAcrossAnswer(grid, clue.Start.Row, clue.Start.Col, clue.Length)
		if gridAnswer != clue.Answer {
			errors = append(errors,
				"across "+string(rune('0'+clue.Number))+": expected "+clue.Answer+", grid has "+gridAnswer)
		}
	}

	for _, clue := range clues.Down {
		gridAnswer := extractDownAnswer(grid, clue.Start.Row, clue.Start.Col, clue.Length)
		if gridAnswer != clue.Answer {
			errors = append(errors,
				"down "+string(rune('0'+clue.Number))+": expected "+clue.Answer+", grid has "+gridAnswer)
		}
	}

	return errors
}

func extractAcrossAnswer(grid [][]Cell, row, col, length int) string {
	var answer strings.Builder
	for i := 0; i < length && col+i < len(grid[0]); i++ {
		answer.WriteString(grid[row][col+i].Solution)
	}
	return answer.String()
}

func extractDownAnswer(grid [][]Cell, row, col, length int) string {
	var answer strings.Builder
	for i := 0; i < length && row+i < len(grid); i++ {
		answer.WriteString(grid[row+i][col].Solution)
	}
	return answer.String()
}
