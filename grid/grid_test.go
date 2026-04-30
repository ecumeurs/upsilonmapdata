package grid

import (
	"testing"

	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

func TestGridBoundaries(t *testing.T) {
	// Create a small grid with height 2 (valid Z: 0, 1)
	g := &Grid{
		Width:  5,
		Length: 5,
		Height: 2,
		Cells:  make(map[position.Position]*cell.Cell),
	}

	// Add a ground cell at the limit (z=1)
	pos1 := position.New(2, 2, 1)
	g.Cells[pos1] = cell.NewCell(cell.Ground, pos1)

	// Add a "leaked" cell at z=2 (outside Height boundary)
	pos2 := position.New(2, 2, 2)
	g.Cells[pos2] = cell.NewCell(cell.Ground, pos2)

	t.Run("PositionIsInGrid", func(t *testing.T) {
		if !g.PositionIsInGrid(pos1) {
			t.Errorf("expected %v to be in grid", pos1)
		}
		if g.PositionIsInGrid(pos2) {
			t.Errorf("expected %v to be out of grid (Height=2)", pos2)
		}
	})

	t.Run("TopMostCellAt respects Height", func(t *testing.T) {
		z := g.TopMostCellAt(2, 2)
		if z != 1 {
			t.Errorf("expected TopMostCellAt to return 1, got %d (should ignore cell at z=2)", z)
		}
	})

	t.Run("CellAt respects Height", func(t *testing.T) {
		_, ok1 := g.CellAt(pos1)
		if !ok1 {
			t.Errorf("expected CellAt to find cell at %v", pos1)
		}

		_, ok2 := g.CellAt(pos2)
		if ok2 {
			t.Errorf("expected CellAt to NOT find cell at %v (out of bounds)", pos2)
		}
	})

	t.Run("MoveEntity respects Height", func(t *testing.T) {
		entityID := uuid.New()
		from := position.New(0, 0, 0)
		g.Cells[from] = cell.NewCell(cell.Ground, from)

		// Move to valid top (z=1)
		err := g.MoveEntity(from, pos1, entityID)
		if err != nil {
			t.Errorf("failed to move to valid position %v: %v", pos1, err)
		}

		// Try to move to "leaked" top (z=2)
		err = g.MoveEntity(pos1, pos2, entityID)
		if err == nil {
			t.Errorf("expected error when moving to out-of-bounds position %v", pos2)
		}
	})
}

func TestGridLowestTopMostConsistency(t *testing.T) {
	g := &Grid{
		Width:  5,
		Length: 5,
		Height: 5,
		Cells:  make(map[position.Position]*cell.Cell),
	}

	// Column with multiple cells
	g.Cells[position.New(0, 0, 0)] = cell.NewCell(cell.Dirt, position.New(0, 0, 0))
	g.Cells[position.New(0, 0, 1)] = cell.NewCell(cell.Dirt, position.New(0, 0, 1))
	g.Cells[position.New(0, 0, 2)] = cell.NewCell(cell.Ground, position.New(0, 0, 2))

	if top := g.TopMostCellAt(0, 0); top != 2 {
		t.Errorf("expected top 2, got %d", top)
	}
	if low := g.LowestCellAt(0, 0); low != 0 {
		t.Errorf("expected low 0, got %d", low)
	}
	if all := g.AllCellsAt(0, 0); len(all) != 3 {
		t.Errorf("expected 3 cells, got %v", all)
	}

	// Column with "leaked" cell at z=5 (Height=5)
	g.Cells[position.New(1, 1, 5)] = cell.NewCell(cell.Ground, position.New(1, 1, 5))
	if top := g.TopMostCellAt(1, 1); top != -1 {
		t.Errorf("expected top -1 (ignoring z=5), got %d", top)
	}
}
