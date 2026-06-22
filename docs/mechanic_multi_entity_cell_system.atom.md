---
id: mechanic_multi_entity_cell_system
human_name: "Multi Entity Cell System"
status: DRAFT
version: 2.0
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
tags: [grid, entities, multi-entity]
parents:
  - [[entity_grid]]
---

# Multi-Entity Cell System

## INTENT

To implement multi-entity cell system where multiple entities (characters, effects, traps, obstacles) can occupy the same grid cell simultaneously, enabling complex tactical interactions.

## THE RULE / LOGIC
**Multi-Entity Cell System Mechanic:**

**Core Principle:**
This mechanic shifts the grid from a "one-entity-per-cell" model to a "multi-occupant" model. Characters, environmental effects, traps, and obstacles can all occupy the same spatial coordinate while maintaining distinct collision and interaction rules.

**Cell Architecture (V2):**
- **Entity Identification:** Each grid cell maintains an `EntityIDs` collection (an array of unique identifiers) rather than a single occupant reference.
- **Effect Identification:** A separate `EffectIDs` collection tracks positional status effects anchored to the cell.
- **Terrain Dynamics:** Cells possess a `CrossingCost` property, allowing varied terrain types (roads, rubble, mud) to influence movement efficiency independently of any active effects.

**Occupancy Constraints:**
- **Character Exclusive Slot:** Each cell supports a maximum of **one** Character entity. Moving a character into a cell already occupied by another character results in a collision failure.
- **Effect/Temporary Stack:** Each cell can hold **multiple** Effect or Temporary entities (e.g. a trap, a fog cloud, and a channeling marker can coexist).
- **Hybrid Occupation:** A cell may simultaneously contain one Character and multiple Temporary entities.

**Collision and Traversal Rules — the WalkThrough property:**
Every entity has a boolean `WalkThrough` attribute governing its physical presence:
- **WalkThrough = True:** allows other entities to pass through / occupy its cell (beneficial clouds, hidden traps, channeling indicators).
- **WalkThrough = False:** acts as a solid obstacle, blocking movement into the cell (Characters, summoned barriers, static walls).

**Entity-default WalkThrough table:**
| Entity type | WalkThrough default | Notes |
|-------------|---------------------|-------|
| Characters | False | always occupy the exclusive character slot |
| Traps | True | must allow entry to trigger the "Step-In" event |
| Environmental Hazards | True or False | True for Poison Gas, False for Ice Wall |
| Markers/Indicators | True | purely informational / pending actions |

**Movement Validation Hierarchy:**
1. **Character Check:** if the cell contains a character other than the mover, movement is denied.
2. **Obstacle (WalkThrough) Check:** iterate all non-character entities; if any has `WalkThrough = False`, movement is denied.
3. **Resolution:** if neither check fails, movement into the cell is permitted.

**Pathfinding Integration (A*):** the navigation algorithm treats a cell as "impassable" if any current occupant is marked non-traversable.

**Cumulative Movement Cost Calculation:**
Total cost to enter a cell = sum of:
1. **Base Step Cost** (typically 1).
2. **Terrain Crossing Cost** (the cell's `CrossingCost`).
3. **Active Effect Penalties** (sum of all `MvtCost` properties from effects in `EffectIDs`).
*Note: a character may enter a cell even if its cost exceeds remaining movement, effectively ending movement there.*

**Grid Lifecycle Management:**
- **Atomic Movement:** `MoveEntity` removes an ID from its source cell's list and appends it to the destination cell's list, keeping the total entity count consistent.
- **Occupancy Retrieval:** specialized lookups return the full collection of entity IDs at any coordinate, enabling area-of-effect and proximity checks.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_multi_entity_cell_system]]`
- **Related Files:** `upsilonmapdata/grid/cell/cell.go`, `upsilonmapdata/grid/grid.go`, `upsilonbattle/battlearena/ruler/rules/move.go`

## EXPECTATION
- A cell accepts at most one Character; moving a character onto a cell already holding another character fails with a collision error.
- A cell may simultaneously hold one Character plus multiple Effect/Temporary entities.
- Movement into a cell is denied if any non-character occupant has `WalkThrough = False`; permitted if all occupants are `WalkThrough = True`.
- Default WalkThrough by type: Characters = False, Traps = True, Markers = True, Hazards = True or False per hazard.
- The cost to enter a cell equals Base Step Cost + cell `CrossingCost` + sum of effect `MvtCost`; a character may still enter when this exceeds remaining movement (ending movement there).
- `MoveEntity` preserves total entity count: the ID is removed from the source cell and appears in exactly the destination cell.
