// Package fill provides a constraint-based crossword grid filler.
package fill

import (
	"lesmotsdatche/internal/domain"
)

// Slot represents a slot in the grid that needs to be filled.
type Slot struct {
	ID        int              // Unique identifier
	Direction domain.Direction // Across or down
	Start     domain.Position  // Starting position
	Length    int              // Number of cells
	Cells     []domain.Position // All cell positions
	Crossings []Crossing       // Intersections with other slots
}

// Crossing represents an intersection between two slots.
type Crossing struct {
	SlotID    int // The other slot's ID
	ThisIndex int // Index in this slot (0-based)
	ThatIndex int // Index in the other slot (0-based)
}

// Pattern returns the current pattern for this slot from the grid.
// Unknown cells are represented as '.'.
func (s *Slot) Pattern(grid [][]rune) string {
	pattern := make([]rune, s.Length)
	for i, pos := range s.Cells {
		c := grid[pos.Row][pos.Col]
		if c == 0 || c == '.' {
			pattern[i] = '.'
		} else {
			pattern[i] = c
		}
	}
	return string(pattern)
}

// IsFilled returns true if all cells in the slot have letters.
func (s *Slot) IsFilled(grid [][]rune) bool {
	for _, pos := range s.Cells {
		c := grid[pos.Row][pos.Col]
		if c == 0 || c == '.' {
			return false
		}
	}
	return true
}

// DiscoverSlots finds all slots in a grid template.
// The grid should have blocks marked and letter cells empty or with existing letters.
func DiscoverSlots(grid [][]domain.Cell) []Slot {
	if len(grid) == 0 {
		return nil
	}

	rows := len(grid)
	cols := len(grid[0])
	var slots []Slot
	slotID := 0

	// Find across slots
	for row := 0; row < rows; row++ {
		col := 0
		for col < cols {
			// Skip blocks
			if grid[row][col].IsBlock() {
				col++
				continue
			}

			// Find start of across slot
			startCol := col
			var cells []domain.Position

			// Extend until block or edge
			for col < cols && grid[row][col].IsLetter() {
				cells = append(cells, domain.Position{Row: row, Col: col})
				col++
			}

			// Only create slot if length >= 2
			if len(cells) >= 2 {
				slots = append(slots, Slot{
					ID:        slotID,
					Direction: domain.DirectionAcross,
					Start:     domain.Position{Row: row, Col: startCol},
					Length:    len(cells),
					Cells:     cells,
				})
				slotID++
			}
		}
	}

	// Find down slots
	for col := 0; col < cols; col++ {
		row := 0
		for row < rows {
			// Skip blocks
			if grid[row][col].IsBlock() {
				row++
				continue
			}

			// Find start of down slot
			startRow := row
			var cells []domain.Position

			// Extend until block or edge
			for row < rows && grid[row][col].IsLetter() {
				cells = append(cells, domain.Position{Row: row, Col: col})
				row++
			}

			// Only create slot if length >= 2
			if len(cells) >= 2 {
				slots = append(slots, Slot{
					ID:        slotID,
					Direction: domain.DirectionDown,
					Start:     domain.Position{Row: startRow, Col: col},
					Length:    len(cells),
					Cells:     cells,
				})
				slotID++
			}
		}
	}

	// Find crossings between slots
	findCrossings(slots)

	return slots
}

// findCrossings populates the Crossings field for each slot.
func findCrossings(slots []Slot) {
	type slotIndex struct {
		slotID int
		index  int
	}

	// Build position to slot index map
	posToSlot := make(map[domain.Position][]slotIndex)

	for i := range slots {
		for j, pos := range slots[i].Cells {
			posToSlot[pos] = append(posToSlot[pos], slotIndex{slotID: i, index: j})
		}
	}

	// Find crossings (positions with 2 slots)
	for _, indices := range posToSlot {
		if len(indices) == 2 {
			s1, s2 := indices[0], indices[1]

			// Add crossing to first slot
			slots[s1.slotID].Crossings = append(slots[s1.slotID].Crossings, Crossing{
				SlotID:    s2.slotID,
				ThisIndex: s1.index,
				ThatIndex: s2.index,
			})

			// Add crossing to second slot
			slots[s2.slotID].Crossings = append(slots[s2.slotID].Crossings, Crossing{
				SlotID:    s1.slotID,
				ThisIndex: s2.index,
				ThatIndex: s1.index,
			})
		}
	}
}

// SlotsByConstraint returns slots sorted by most constrained first.
// More constrained = fewer possible candidates = should be filled first.
func SlotsByConstraint(slots []Slot, grid [][]rune, lexicon Lexicon) []int {
	type slotScore struct {
		id    int
		score int // Lower = more constrained
	}

	scores := make([]slotScore, 0, len(slots))

	for i, slot := range slots {
		if slot.IsFilled(grid) {
			continue // Skip already filled slots
		}

		pattern := slot.Pattern(grid)
		candidates := lexicon.Match(pattern)
		scores = append(scores, slotScore{id: i, score: len(candidates)})
	}

	// Sort by score (ascending = most constrained first)
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].score < scores[i].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	result := make([]int, len(scores))
	for i, s := range scores {
		result[i] = s.id
	}
	return result
}
