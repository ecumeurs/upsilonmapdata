---
id: entity_mapdata_3d_grid
status: DRAFT
dependents: []
human_name: "MapData 3D Grid Entity"
version: 1.0
tags: [mapdata,grid,entity]
parents:
  - [[rule_mapdata_grid_standard]]
type: ENTITY
layer: ARCHITECTURE
priority: 2
---

# MapData 3D Grid Entity

## INTENT
Provide a concrete data structure for representing the tactical grid in memory.

## THE RULE / LOGIC
- **Structure:**
  - `Width`, `Length`, `Height` (int)
  - `Cells` (Map[Position]Cell)
- **Validation:** Coordinates must be within [0, Width-1], [0, Length-1], [0, Height-1].
- **Serialization:** Must be compatible with JSON for bridge transmission.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[entity_mapdata_3d_grid]]`

## EXPECTATION
The Grid entity must successfully store and retrieve cells by position and provide helper methods for spatial queries (TopMost, Lowest).
