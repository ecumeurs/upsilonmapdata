---
id: rule_mapdata_grid_standard
status: DRAFT
type: RULE
version: 1.0
priority: 2
tags: [mapdata,grid,standard]
dependents:
  - [[entity_mapdata_3d_grid]]
human_name: "MapData Grid Standard"
layer: BUSINESS
parents:
  - [[contract_mapdata_contract]]
---

# MapData Grid Standard

## INTENT
Define the mandatory structural properties for tactical grids in Upsilon.

## THE RULE / LOGIC
- **Dimensions:** Every grid must have explicit Width (X), Length (Y), and Height (Z) limits.
- **Cell Mapping:** Coordinates must be mappable to unique cell instances.
- **Verticality:** The grid must support multiple vertical layers (Z > 1) to represent diverse terrain.
- **Obstacles:** The standard must support a way to distinguish walkable tiles from impassable obstacles.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[rule_mapdata_grid_standard]]`

## EXPECTATION
Grid structures must consistently support Width, Length, and Height properties and validate all cell additions against these boundaries.
