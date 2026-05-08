// Package grid provides pathfinding and spatial selection logic for the tactical grid.
// It implements A* search and pattern-based cell filtering.
// @spec-link [[mapdata_grid_standard]]
// @spec-link [[mapdata_3d_grid]]
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
	// Apply the pattern and filter by grid dimensions.
	pos := pat.ApplyInArea(origin, g.Width, g.Length, g.Height)
	for _, p := range pos {
		// Only include positions that actually have a cell associated with them.
		if g.Contains(p) {
			res = append(res, p)
		}
	}
	return res
}

// SelectPositionsByPattern2D returns the top-most ground positions matching the 2D pattern.
// This is used for skill targeting where the user selects a (X,Y) tile and the engine finds the surface.
func (g *Grid) SelectPositionsByPattern2D(origin position.Position, pat pattern.Pattern2D) []position.Position {
	res := []position.Position{}
	for _, p := range pat {
		pos := origin.Add(p)
		// Snap the position to the top-most cell at those horizontal coordinates.
		pos.Z = g.TopMostCellAt(pos.X, pos.Y)
		if g.Contains(pos) {
			res = append(res, pos)
		}
	}
	return res
}

// CellsForPositions returns the cell data for a list of positions.
// It safely handles cases where some positions might not correspond to a valid cell.
func (g *Grid) CellsForPositions(pos []position.Position) []*cell.Cell {
	res := []*cell.Cell{}
	for _, p := range pos {
		// Lookup the cell in the master map.
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
	// 1. Initial sanity checks.
	if !g.Contains(start) || !g.Contains(end) {
		return nil, false
	}
	if start == end {
		return []position.Position{start}, true
	}

	// 2. Initialize search state.
	visited := map[position.Position]int{}
	queue := []position.Position{start}
	parents := map[position.Position]position.Position{}
	
	// 3. Main search loop.
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		
		// If we reached the goal, reconstruct and return the path.
		if pos.Equals(end) {
			return g.reconstructPath(visited, start, end, jumpHeight, exclude), true
		}
		
		// Mark current node as visited and explore neighbors.
		visited[pos] = visited[parents[pos]] + 1
		queue = g.exploreAStarNeighbors(pos, end, jumpHeight, exclude, visited, parents, queue)
	}
	return nil, false
}

// exploreAStarNeighbors is a helper to process adjacent tiles during pathfinding.
func (g *Grid) exploreAStarNeighbors(pos, end position.Position, jump int, ex func(position.Position) bool, vis map[position.Position]int, pars map[position.Position]position.Position, q []position.Position) []position.Position {
	for _, n := range g.SelectPositionsByPattern2D(pos, pattern.Neighbours2D()) {
		// Filter out invalid or already visited neighbors.
		if _, found := vis[n]; found {
			continue
		}
		if tools.Abs(n.Z-pos.Z) > jump {
			continue
		}
		if ex != nil && ex(n) && !n.Equals(end) {
			continue
		}
		
		// Only walkable ground tiles are valid for standard pathfinding.
		if c, found := g.CellAt(n); found && c.Type == cell.Ground {
			q = append(q, n)
			pars[n] = pos
		}
	}
	return q
}

// reconstructPath builds the full path sequence from the A* visited nodes.
func (g *Grid) reconstructPath(visited map[position.Position]int, start, end position.Position, jumpHeight int, exclude func(position.Position) bool) []position.Position {
	res := []position.Position{end}
	// Trace back from the end position using the adjacency and distance metrics.
	for res[len(res)-1] != start {
		lowestPos := g.findLowestAdjacent(res[len(res)-1], jumpHeight, visited)
		res = append(res, lowestPos)
	}

	// Reverse the result to get the path from start to end.
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

// findLowestAdjacent is a helper for path reconstruction to find the neighbor closest to the start.
func (g *Grid) findLowestAdjacent(curr position.Position, jump int, vis map[position.Position]int) position.Position {
	lowest := 999999
	var lowestPos position.Position
	for _, n := range g.SelectPositionsByPattern2D(curr, pattern.Neighbours2D()) {
		if tools.Abs(n.Z-curr.Z) <= jump {
			if v, found := vis[n]; found && v < lowest {
				lowest = v
				lowestPos = n
			}
		}
	}
	return lowestPos
}

