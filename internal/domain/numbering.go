package domain

// AssignNumbers assigns clue numbers to a grid in row-major order.
// A cell gets a number if it starts an across entry, a down entry, or both.
//
// A cell starts an across entry if:
//   - It is a letter cell
//   - Its left neighbor is a block or out of bounds
//   - It has at least one letter cell to its right
//
// A cell starts a down entry if:
//   - It is a letter cell
//   - Its top neighbor is a block or out of bounds
//   - It has at least one letter cell below it
//
// The function returns a new grid with Number fields populated.
// The original grid is not modified.
func AssignNumbers(grid [][]Cell) [][]Cell {
	if len(grid) == 0 {
		return nil
	}

	rows := len(grid)
	cols := len(grid[0])

	// Create a deep copy of the grid
	result := make([][]Cell, rows)
	for i := range grid {
		result[i] = make([]Cell, cols)
		copy(result[i], grid[i])
		// Clear any existing numbers
		for j := range result[i] {
			result[i][j].Number = 0
		}
	}

	currentNumber := 1

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cell := &result[row][col]

			// Skip blocks
			if cell.IsBlock() {
				continue
			}

			startsAcross := startsAcrossEntry(result, row, col, rows, cols)
			startsDown := startsDownEntry(result, row, col, rows, cols)

			if startsAcross || startsDown {
				cell.Number = currentNumber
				currentNumber++
			}
		}
	}

	return result
}

// startsAcrossEntry checks if a cell starts an across entry.
func startsAcrossEntry(grid [][]Cell, row, col, rows, cols int) bool {
	// Must be a letter cell
	if grid[row][col].IsBlock() {
		return false
	}

	// Left must be block or out of bounds
	leftIsBarrier := col == 0 || grid[row][col-1].IsBlock()
	if !leftIsBarrier {
		return false
	}

	// Must have at least one letter cell to the right
	hasRightLetter := col+1 < cols && grid[row][col+1].IsLetter()
	return hasRightLetter
}

// startsDownEntry checks if a cell starts a down entry.
func startsDownEntry(grid [][]Cell, row, col, rows, cols int) bool {
	// Must be a letter cell
	if grid[row][col].IsBlock() {
		return false
	}

	// Top must be block or out of bounds
	topIsBarrier := row == 0 || grid[row-1][col].IsBlock()
	if !topIsBarrier {
		return false
	}

	// Must have at least one letter cell below
	hasBelowLetter := row+1 < rows && grid[row+1][col].IsLetter()
	return hasBelowLetter
}

// StartsAcross returns true if the cell at (row, col) starts an across entry.
func StartsAcross(grid [][]Cell, row, col int) bool {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return false
	}
	return startsAcrossEntry(grid, row, col, len(grid), len(grid[0]))
}

// StartsDown returns true if the cell at (row, col) starts a down entry.
func StartsDown(grid [][]Cell, row, col int) bool {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return false
	}
	return startsDownEntry(grid, row, col, len(grid), len(grid[0]))
}
