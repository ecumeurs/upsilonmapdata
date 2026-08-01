# upsilonmapdata

A Go library providing the 3D tactical grid data structures (cells, positions, and area/line patterns) used by the Upsilon battle engine.

## What it does

`upsilonmapdata` models a battlefield as a 3D `Grid` of `Cell`s indexed by `(X, Y, Z)` positions:

- **`grid`** — the `Grid` type itself: bootstrapping a flat map with ground/dirt layers, looking up cells by position, moving/adding/removing entities per cell, A* pathfinding with jump-height limits, and an `obelisk.js`-based HTML renderer for debugging/previewing a map.
- **`grid/cell`** — the `Cell` type: terrain type (`Obstacle`, `Ground`, `Water`, `Dirt`, `Debug`, `Debug2`), and the entities/effects currently occupying it.
- **`grid/position`** — the `Position` type: a discrete 3D coordinate with arithmetic (add/subtract), distance, adjacency, and random-position helpers.
- **`grid/position/pattern`** — `Pattern` (3D) and `Pattern2D` (flat) shapes — circles, squares, lines, neighbour sets — used to compute area-of-effect selections, targeting, and terrain generation relative to an origin.

It has no server, no persistence, and no game rules of its own — it's a shared data/geometry layer consumed by other Upsilon services (e.g. the battle engine and map tooling) to represent and query the state of a map.

## Usage

```go
import (
    "github.com/ecumeurs/upsilonmapdata/grid"
    "github.com/ecumeurs/upsilonmapdata/grid/position"
    "github.com/ecumeurs/upsilonmapdata/grid/position/pattern"
)

// Build a 10x10 flat grid with ground at level 1.
g := grid.NewGrid(10, 10, 1)

// Query and mutate cells.
p := position.New(3, 4, 1)
if c, ok := g.CellAt(p); ok {
    c.AddEntity(entityID)
}

// Select tiles in a 2-tile radius around a position.
targets := g.SelectPositionsByPattern(p, pattern.Circle(2))

// Find a path between two positions.
path, found := g.AStarPath(start, end, jumpHeight, nil)
```

Import the module as `github.com/ecumeurs/upsilonmapdata` (Go 1.25+).
