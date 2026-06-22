// Package grid provides automated tests for the tactical grid system.
// @spec-link [[rule_mapdata_grid_standard]]
// @spec-link [[entity_mapdata_3d_grid]]
package grid

import (
	"testing"

	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// TestGridBoundaries verifies that spatial queries correctly respect the grid's hard dimensions.
// It specifically checks the edge cases for Z-axis height limits and boundary enforcement.
// This test ensures that the grid behaves as a finite volume rather than an infinite map.
func TestGridBoundaries(t *testing.T) {
	// 1. Setup Phase: Create a small grid with height 2 (valid Z range: 0, 1).
	// We use a small volume to make edge cases easier to trigger and verify.
	g := &Grid{
		Width:  5,
		Length: 5,
		Height: 2,
		Cells:  make(map[position.Position]*cell.Cell),
	}

	// 2. Data Insertion: Add a ground cell at the limit (z=1).
	pos1 := position.New(2, 2, 1)
	g.Cells[pos1] = cell.NewCell(cell.Ground, pos1)

	// 3. Negative Case: Add a "leaked" cell at z=2 (outside Height boundary).
	// This simulates a potential data inconsistency that the logic must handle gracefully.
	pos2 := position.New(2, 2, 2)
	g.Cells[pos2] = cell.NewCell(cell.Ground, pos2)

	// 4. PositionIsInGrid Validation:
	// Verify that coordinates are correctly identified as being within or outside the finite grid.
	t.Run("PositionIsInGrid", func(t *testing.T) {
		// Valid position (z=1) must return true.
		if !g.PositionIsInGrid(pos1) {
			t.Errorf("expected %v to be in grid", pos1)
		}
		// Invalid position (z=2) must return false even if a cell exists there in memory.
		if g.PositionIsInGrid(pos2) {
			t.Errorf("expected %v to be out of grid (Height=2)", pos2)
		}
	})

	// 5. TopMostCellAt Logic Check:
	// Verify that vertical scanning correctly stops at the defined grid height limit.
	t.Run("TopMostCellAt respects Height", func(t *testing.T) {
		z := g.TopMostCellAt(2, 2)
		// The cell at z=2 should be ignored because the grid height is capped at 2.
		if z != 1 {
			t.Errorf("expected TopMostCellAt to return 1, got %d (should ignore cell at z=2)", z)
		}
	})

	// 6. CellAt Lookup Verification:
	// Ensure that high-level lookup methods enforce boundary rules on the coordinate level.
	t.Run("CellAt respects Height", func(t *testing.T) {
		// Valid coordinates must yield the cell data.
		_, ok1 := g.CellAt(pos1)
		if !ok1 {
			t.Errorf("expected CellAt to find cell at %v", pos1)
		}

		// Coordinates outside [0, Height) must return false/nil even if the map entry exists.
		_, ok2 := g.CellAt(pos2)
		if ok2 {
			t.Errorf("expected CellAt to NOT find cell at %v (out of bounds)", pos2)
		}
	})

	// 7. Entity Movement Enforcement:
	// Verify that the engine prevents actors from moving into non-playable spatial sectors.
	t.Run("MoveEntity respects Height", func(t *testing.T) {
		entityID := uuid.New()
		from := position.New(0, 0, 0)
		g.Cells[from] = cell.NewCell(cell.Ground, from)

		// Successful Move: Transition from origin to the highest valid tile (z=1).
		err := g.MoveEntity(from, pos1, entityID)
		if err != nil {
			t.Errorf("failed to move to valid position %v: %v", pos1, err)
		}

		// Blocked Move: Attempt to move higher into an out-of-bounds coordinate (z=2).
		err = g.MoveEntity(pos1, pos2, entityID)
		if err == nil {
			t.Errorf("expected error when moving to out-of-bounds position %v", pos2)
		}
	})
	// 8. Conclusion: Boundary test sequence successfully verified all limits.
	// This ensures that spatial queries are safe and predictable across the grid.
}


// TestGridLowestTopMostConsistency ensures that queries for the top/bottom cells of a column are consistent.
// It validates that the grid's vertical traversal algorithms correctly identify stack boundaries.
func TestGridLowestTopMostConsistency(t *testing.T) {
	// 1. Grid Initialization: Create a tall grid column for testing stacking behavior.
	g := &Grid{
		Width:  5,
		Length: 5,
		Height: 5,
		Cells:  make(map[position.Position]*cell.Cell),
	}

	// 2. Stack Creation: Insert three cells into the (0, 0) column at different heights.
	// This represents a multi-layered geological stack (e.g., dirt with ground on top).
	g.Cells[position.New(0, 0, 0)] = cell.NewCell(cell.Dirt, position.New(0, 0, 0))
	g.Cells[position.New(0, 0, 1)] = cell.NewCell(cell.Dirt, position.New(0, 0, 1))
	g.Cells[position.New(0, 0, 2)] = cell.NewCell(cell.Ground, position.New(0, 0, 2))

	// 3. Highest Cell Check:
	// Verify that the TopMostCellAt algorithm correctly identifies the z=2 cell.
	if top := g.TopMostCellAt(0, 0); top != 2 {
		t.Errorf("expected top 2, got %d", top)
	}
	// 4. Lowest Cell Check:
	// Verify that the LowestCellAt algorithm correctly identifies the base cell at z=0.
	if low := g.LowestCellAt(0, 0); low != 0 {
		t.Errorf("expected low 0, got %d", low)
	}
	// 5. Column Density Check:
	// Ensure that AllCellsAt returns the correct count of stacked objects in the vertical column.
	if all := g.AllCellsAt(0, 0); len(all) != 3 {
		t.Errorf("expected 3 cells, got %v", all)
	}

	// 6. Exclusion Check:
	// Verify that cells outside the playable Height (z=5 in a Height=5 grid) are ignored.
	// This hardens the system against floating-point or off-by-one errors in map data.
	g.Cells[position.New(1, 1, 5)] = cell.NewCell(cell.Ground, position.New(1, 1, 5))
	if top := g.TopMostCellAt(1, 1); top != -1 {
		t.Errorf("expected top -1 (ignoring z=5), got %d", top)
	}
	// Consistency verification complete.
}


