package grid

import (
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilonmapdata/grid/position/pattern"
	"github.com/ecumeurs/upsilontools/tools"
)

// SelectPositionsByPattern returns the positions matching the pattern within the grid boundaries.
func (g *Grid) SelectPositionsByPattern(origin position.Position, pat pattern.Pattern) []position.Position {
	res := []position.Position{}
	pos := pat.ApplyInArea(origin, g.Width, g.Length, g.Height)
	for _, p := range pos {
		if g.Contains(p) {
			res = append(res, p)
		}
	}
	return res
}

// SelectPositionsByPattern2D returns the top-most ground positions matching the 2D pattern.
func (g *Grid) SelectPositionsByPattern2D(origin position.Position, pat pattern.Pattern2D) []position.Position {
	res := []position.Position{}
	for _, p := range pat {
		pos := origin.Add(p)
		pos.Z = g.TopMostCellAt(pos.X, pos.Y)
		if g.Contains(pos) {
			res = append(res, pos)
		}
	}
	return res
}

// CellsForPositions returns the cell data for a list of positions.
func (g *Grid) CellsForPositions(pos []position.Position) []*cell.Cell {
	res := []*cell.Cell{}
	for _, p := range pos {
		c, ok := g.Cells[p]
		if ok {
			res = append(res, c)
		}
	}
	return res
}

// AStarPath returns the shortest path from start to end using the A* algorithm.
// It considers jump height limits and optional exclusion criteria for blocked tiles.
func (g *Grid) AStarPath(start, end position.Position, jumpHeight int, exclude func(position.Position) bool) ([]position.Position, bool) {
	if !g.Contains(start) || !g.Contains(end) {
		return nil, false
	}
	if start == end {
		return []position.Position{start}, true
	}

	visited := map[position.Position]int{}
	queue := []position.Position{start}
	parents := map[position.Position]position.Position{}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		if pos.Equals(end) {
			return g.reconstructPath(visited, start, end, jumpHeight, exclude), true
		}
		visited[pos] = visited[parents[pos]] + 1
		for _, n := range g.SelectPositionsByPattern2D(pos, pattern.Neighbours2D()) {
			if _, found := visited[n]; found {
				continue
			}
			if tools.Abs(n.Z-pos.Z) > jumpHeight {
				continue
			}
			// Skip if it matches the exclusion criteria (unless it's the target)
			if exclude != nil && exclude(n) && !n.Equals(end) {
				continue
			}
			if c, found := g.CellAt(n); found && c.Type == cell.Ground {
				queue = append(queue, n)
				parents[n] = pos
			}
		}
	}
	return nil, false
}

// reconstructPath builds the full path sequence from the A* visited nodes.
func (g *Grid) reconstructPath(visited map[position.Position]int, start, end position.Position, jumpHeight int, exclude func(position.Position) bool) []position.Position {
	res := []position.Position{end}
	for res[len(res)-1] != start {
		// find the lowest number within adjascents
		lowest := 999999
		var lowestPos position.Position
		for _, n := range g.SelectPositionsByPattern2D(res[len(res)-1], pattern.Neighbours2D()) {
			if tools.Abs(n.Z-res[len(res)-1].Z) <= jumpHeight {
				if vis, found := visited[n]; found && vis < lowest {
					lowest = vis
					lowestPos = n
				}
			}
		}
		res = append(res, lowestPos)
	}

	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
