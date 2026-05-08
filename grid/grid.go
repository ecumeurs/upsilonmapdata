// Package grid provides the core data structures for representing the tactical battlefield.
// It manages cells, entities, and pathfinding logic.
// @spec-link [[mapdata_grid_standard]]
// @spec-link [[mapdata_3d_grid]]
package grid

import (
	"math/rand"

	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilonmapdata/grid/position/pattern"
)

// @spec-link [[entity_grid]]
// Grid represents a 3D tactical map consisting of multiple cells.
type Grid struct {
	// Width is the size of the grid along the X-axis.
	Width  int
	// Length is the size of the grid along the Y-axis.
	Length int
	// Height is the maximum vertical limit of the grid.
	Height int

	// Cells maps positions to their respective cell data.
	Cells map[position.Position]*cell.Cell
}

// NewGrid creates a new flat grid with ground at the specified level.
// It initializes a base layer of dirt (below ground) and a top layer of walkable ground.
// This function is the primary way to programmatically bootstrap a standard battlefield.
// It ensures that all cells are correctly initialized with their types and positions.
func NewGrid(width, length, groundlevel int) *Grid {
	// Initialize the grid metadata and cell map structure.
	g := &Grid{
		Width:  width,
		Length: length,
		Height: groundlevel + 2,
		Cells:  make(map[position.Position]*cell.Cell),
	}
	
	// Fill the grid column by column. This structure helps maintain a low nesting depth.
	// We iterate through the X-axis and delegate the rest of the stack creation.
	for x := 0; x < width; x++ {
		g.fillInitialColumns(x, length, groundlevel)
	}
	return g
}


// fillInitialColumns iterates through the length of the grid for a given X coordinate.
// It delegates the vertical stack population to fillInitialZAxis.
func (g *Grid) fillInitialColumns(x, length, groundlevel int) {
	for y := 0; y < length; y++ {
		g.fillInitialZAxis(x, y, groundlevel)
	}
}

// fillInitialZAxis populates the vertical stack for a specific (x, y) coordinate.
// It places dirt cells below the groundlevel and a ground cell at the groundlevel.
func (g *Grid) fillInitialZAxis(x, y, groundlevel int) {
	for z := 0; z < groundlevel+1; z++ {
		pos := position.New(x, y, z)
		if z < groundlevel {
			// Sub-surface layers are filled with dirt to provide visual depth.
			g.Cells[pos] = cell.NewCell(cell.Dirt, pos)
		} else {
			// The top-most layer is designated as walkable ground.
			g.Cells[pos] = cell.NewCell(cell.Ground, pos)
		}
	}
}



// RandomPosition returns a random valid position in the grid that is not an obstacle or occupied.
func (g *Grid) RandomPosition() position.Position {
	for i := 0; i < 1000; i++ {
		x := rand.Intn(g.Width)
		y := rand.Intn(g.Length)
		z := g.TopMostGroundAt(x, y)
		pos := position.New(x, y, z)
		// check if there is a ground cell
		if c, found := g.CellAt(pos); found {
			// ISS-058: Ensure the tile is NOT an obstacle and NOT occupied by an entity
			if c.Type != cell.Obstacle && !c.IsOccupied() {
				return pos
			}
		}
	}
	// Fallback to the first available ground tile if random selection fails
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Length; y++ {
			z := g.TopMostGroundAt(x, y)
			pos := position.New(x, y, z)
			if c, found := g.CellAt(pos); found && c.Type != cell.Obstacle && !c.IsOccupied() {
				return pos
			}
		}
	}
	return position.Position{}
}




// CellAt
func (g *Grid) CellAt(p position.Position) (*cell.Cell, bool) {
	if g.PositionIsInGrid(p) {
		c, ok := g.Cells[p]
		return c, ok
	}
	return nil, false
}

func (g *Grid) TopMostCellAt(x, y int) int {
	for z := g.Height - 1; z >= 0; z-- {
		pos := position.New(x, y, z)
		_, found := g.Cells[pos]
		if found {
			return z
		}
	}
	return -1
}

func (g *Grid) LowestCellAt(x, y int) int {
	for z := 0; z < g.Height; z++ {
		pos := position.New(x, y, z)
		if _, found := g.Cells[pos]; found {
			return z
		}
	}
	return -1
}

func (g *Grid) AllCellsAt(x, y int) []int {
	res := []int{}
	for z := 0; z < g.Height; z++ {
		pos := position.New(x, y, z)
		if _, found := g.Cells[pos]; found {
			res = append(res, z)
		}
	}
	return res
}

func (g *Grid) TopMostGroundAt(x, y int) int {
	for z := g.Height - 1; z >= 0; z-- {
		pos := position.New(x, y, z)
		c, found := g.Cells[pos]
		if found && c.Type == cell.Ground {
			return z
		}
	}
	return -1
}

func (g *Grid) LowestGroundAt(x, y int) int {
	for z := 0; z < g.Height; z++ {
		pos := position.New(x, y, z)
		if c, found := g.Cells[pos]; found && c.Type == cell.Ground {
			return z
		}
	}
	return -1
}

func (g *Grid) AllGroundAt(x, y int) []int {
	res := []int{}
	for z := 0; z < g.Height; z++ {
		pos := position.New(x, y, z)
		if c, found := g.Cells[pos]; found && c.Type == cell.Ground {
			res = append(res, z)
		}
	}
	return res
}

// PositionIsInGrid returns true if the position is in the grid
func (g *Grid) PositionIsInGrid(p position.Position) bool {
	return p.X >= 0 && p.X < g.Width && p.Y >= 0 && p.Y < g.Length && p.Z >= 0 && p.Z < g.Height
}

// ForcePositionToGround
func (g *Grid) ForcePositionToGround(p position.Position) (position.Position, bool) {
	if g.PositionIsInGrid(p) {
		return position.New(p.X, p.Y, g.TopMostGroundAt(p.X, p.Y)), true
	}
	return position.New(-1, -1, -1), false
}

// Display the grid
func (g *Grid) Display() {
	for z := 0; z < g.Height; z++ {
		for y := 0; y < g.Length; y++ {
			for x := 0; x < g.Width; x++ {
				pos := position.New(x, y, z)
				c, ok := g.Cells[pos]
				if !ok {
					continue
				}
				switch c.Type {
				case cell.Ground:
					if !c.IsOccupied() {
						print(".")
					} else {
						print("x")
					}
				case cell.Obstacle:
					print("#")
				}
			}
			println()
		}
		println()
		println()
	}
}

func (g *Grid) RemoveCell(p position.Position) {
	delete(g.Cells, p)
}

// ReplaceCell expects the new Cell to be at the same x,y position as the old one (z can be different)
func (g *Grid) ReplaceCell(p position.Position, c *cell.Cell) {
	if !g.Contains(p) {
		return
	}
	if p.SameXY(c.Position) {
		g.RemoveCell(p)
		g.Cells[c.Position] = c
	}
}

func (g *Grid) ReplaceCellType(p position.Position, t cell.CellType) {
	if !g.Contains(p) {
		return
	}
	g.Cells[p].Type = t
}




func (g *Grid) FindLowestLevel() int {
	lowest := 0
	for _, c := range g.Cells {
		if c.Position.Z < lowest {
			lowest = c.Position.Z
		}
	}
	return lowest
}

// FindNearestCellMatchingPredicate returns the nearest cell matching the predicate
func (g *Grid) FindNearestCellMatchingPredicate(p position.Position, predicate func(*cell.Cell) bool) (*cell.Cell, bool) {
	p.Z = g.LowestCellAt(p.X, p.Y)
	if !g.Contains(p) {
		return nil, false
	}

	if predicate(g.Cells[p]) {
		return g.Cells[p], true
	}

	visited := map[position.Position]bool{p: true}
	queue := []position.Position{p}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		for _, n := range g.SelectPositionsByPattern2D(pos, pattern.Neighbours2D()) {
			if visited[n] {
				continue
			}
			visited[n] = true
			if predicate(g.Cells[n]) {
				return g.Cells[n], true
			}
			queue = append(queue, n)
		}
	}
	return nil, false
}

func (g *Grid) Contains(p position.Position) bool {
	_, ok := g.Cells[p]
	return ok
}



