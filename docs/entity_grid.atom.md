---
id: entity_grid
human_name: Grid Entity
type: ENTITY
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: [spatial, map, 3d]
parents:
  - [[rule_mapdata_grid_standard]]
dependents:
  - [[mechanic_multi_entity_cell_system]]
---
# Grid Entity

## INTENT
To manage a 3D spatial collection of cells, providing utilities for navigation, entity placement, and layout management.

## THE RULE / LOGIC
- **Dimensionality:** A Grid is defined by `Width` (X), `Length` (Y), and `Height` (Z).
- **Cell Collection:** A Grid contains a mapping of `position.Position` (X, Y, Z) to `cell.Cell`.
- **Terrain / Obstacle Distinction:** Every cell carries a `CellType` — `Ground` (walkable), `Obstacle` (impassable, blocks line-of-sight), `Water`, `Dirt` (sub-surface fill), `Debug`/`Debug2` (diagnostic markers). Pathfinding and random-position selection only ever treat `Ground` as traversable/eligible; `Obstacle` is never selected or crossed.
- **Verticality:**
    - `TopMostGroundAt(x, y)` / `LowestGroundAt(x, y)`: Identifies the highest/lowest Z-index where a `Ground` type cell exists.
    - `TopMostCellAt(x, y)` / `LowestCellAt(x, y)`: Identifies the highest/lowest Z-index where *any* cell exists, regardless of type — this is what the 2D client projection and pattern-snapping use, distinct from the ground-only lookups above.
- **Entity Management:**
    - Entities are pinned to specific cells.
    - `MoveEntity(from, to, uuid)`: Updates `EntityID` in the source and destination cells if valid.
- **Pathfinding:**
    - `AStarPath(start, end, jumpHeight, exclude)`: Calculates a path between two positions, restricted by a maximum vertical step (`jumpHeight`), an optional per-position exclusion predicate, and traversable cell types (`Ground` only).
- **Spatial Patterns (2D/3D position selection):** `SelectPositionsByPattern`/`SelectPositionsByPattern2D` resolve a relative `Pattern` (3D offsets: circle, square, line, neighbours, etc.) or `Pattern2D` against an origin into concrete in-grid positions — this is the mechanism skills and area-of-effect calculations use to target multiple tiles. `Pattern2D` additionally snaps every offset to the top-most cell at that (X,Y) column via `TopMostCellAt`, so flat/2D-style targeting always resolves to the visible surface regardless of Z.
- **Boundaries:** All operations must verify if a `position.Position` is within `[0, Width)`, `[0, Length)`, and `[0, Height)`.
- **Visualization (obelisk):** The grid can render itself as a standalone HTML/JS debug view (`GenerateHTML`) using the `obelisk.js` isometric-cube library — one colored cube per cell, color and shading driven by `CellType` and Z-height. This is a debugging/preview aid only; it is not consumed by any production client.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_grid]]`
- **Primary Struct:** `type Grid struct` in `grid.go`
- **Supporting Files:** `grid/cell/cell.go` (`CellType`, `Cell`), `grid/position/position.go` (`Position`, distance/adjacency), `grid/position/pattern/pattern.go` + `pattern_2d.go` (`Pattern`/`Pattern2D` selection), `grid/entity_manager.go` (entity placement/movement), `grid/pathfinding.go` (A* + pattern-based selection), `grid/obelisk.go` (HTML/obelisk.js visualization)
- **Related Atoms:** `[[rule_mapdata_grid_standard]]`, `[[mechanic_multi_entity_cell_system]]`, `[[mechanic_cell_attached_effects]]`

## EXPECTATION (For Testing)
- `NewGrid(10, 10, 2)` -> Returns a grid where `Height` is 4 (groundlevel + 2) and all X,Y positions have `Dirt` or `Ground` cells up to level 2.
- `AStarPath` from (0,0,0) to (1,1,1) with `jumpHeight=0` -> Fails if verticality change is required; the returned path never crosses an `Obstacle` cell.
- `MoveEntity` to a position outside `Width`/`Length` -> Returns an error.
- `SelectPositionsByPattern`/`SelectPositionsByPattern2D` never return a position outside the grid or lacking a backing cell; `Obstacle` cells are excluded from `AStarPath` results and from `RandomPosition` selection.

## CLIENT PROJECTION
Clients (CLI, battleui) consume a 2D projection of this grid: for each `(x, y)` column, the API returns the **topmost cell** (via `TopMostCellAt`) plus its Z index as `height`, and the full engine Z ceiling as `max_height`. Caves/underground are not exposed in the current iteration.
