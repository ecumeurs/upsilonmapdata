// Package cell defines the individual components of the tactical grid.
// Each cell maintains its own type, position, and lists of attached entities and effects.
// @spec-link [[rule_mapdata_grid_standard]]
// @spec-link [[entity_mapdata_3d_grid]]
package cell

import (
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// CellType defines the terrain category of a grid cell.
type CellType int

const (
	// Obstacle tiles are impassable and block line-of-sight.
	Obstacle CellType = 0
	// Ground tiles are walkable and serve as the primary combat layer.
	Ground CellType = 1
	// Water tiles may have movement penalties or unique interactions.
	Water  CellType = 2
	// Dirt tiles represent the sub-surface layers of the grid.
	Dirt   CellType = 3
	// Debug tiles are used for visual diagnostic markers.
	Debug  CellType = 4
	// Debug2 tiles are alternative diagnostic markers.
	Debug2 CellType = 5
)

// Cell represents a single atomic unit of the tactical battlefield.
// It supports multiple entities (characters, traps) and persistent effects (zones).
//
// @spec-link [[mechanic_multi_entity_cell_system]]
// @spec-link [[mechanic_cell_attached_effects]]
type Cell struct {
	// Type determines the collision and visual properties of the cell.
	Type     CellType
	// Position is the 3D coordinate of the cell within the grid.
	Position position.Position

	// EntityIDs holds all entities currently occupying this space.
	// This includes characters, items, and interactive objects.
	EntityIDs []uuid.UUID

	// EffectIDs holds references to area-of-effect or persistent terrain modifiers.
	EffectIDs []uuid.UUID
}

// NewCell initializes a new cell instance with the specified terrain type and 3D position.
// It performs mandatory slice initialization for entities and effects to prevent nil panics.
// This is the primary constructor for generating game-ready grid cells.
func NewCell(t CellType, p position.Position) *Cell {
	// 1. Assign the provided terrain type and world coordinates.
	// 2. Initialize the entity registry as an empty, non-nil slice.
	// 3. Initialize the effect registry as an empty, non-nil slice.
	return &Cell{
		Type:      t,
		Position:  p,
		EntityIDs: []uuid.UUID{},
		EffectIDs: []uuid.UUID{},
	}
}

// AddEntity appends a new entity ID to the cell's registry.
func (c *Cell) AddEntity(id uuid.UUID) {
	c.EntityIDs = append(c.EntityIDs, id)
}

// RemoveEntity deletes an entity ID from the cell's registry if present.
func (c *Cell) RemoveEntity(id uuid.UUID) {
	c.EntityIDs = removeID(c.EntityIDs, id)
}

// HasEntity returns true if the specified entity ID is currently registered in this cell.
func (c *Cell) HasEntity(id uuid.UUID) bool {
	for _, eid := range c.EntityIDs {
		if eid == id {
			return true
		}
	}
	return false
}

// IsOccupied returns true if the cell contains at least one entity.
func (c *Cell) IsOccupied() bool {
	return len(c.EntityIDs) > 0
}

// AddEffect attaches a new positional effect ID to the cell.
func (c *Cell) AddEffect(id uuid.UUID) {
	c.EffectIDs = append(c.EffectIDs, id)
}

// RemoveEffect detaches a positional effect ID from the cell.
func (c *Cell) RemoveEffect(id uuid.UUID) {
	c.EffectIDs = removeID(c.EffectIDs, id)
}

// ToString returns the human-readable string representation of the CellType.
func (ct *CellType) ToString() string {
	switch *ct {
	case Obstacle:
		return "Obstacle"
	case Ground:
		return "Ground"
	case Water:
		return "Water"
	case Dirt:
		return "Dirt"
	case Debug:
		return "Debug"
	case Debug2:
		return "Debug2"
	default:
		return "Unknown"
	}
}

// removeID is an internal helper that removes a target UUID from a slice while preserving relative order.
// It uses slice reslicing to perform the removal in O(n) time.
func removeID(ids []uuid.UUID, id uuid.UUID) []uuid.UUID {
	// 1. Iterate through the slice to find the target identifier.
	for i, existing := range ids {
		// 2. Exact UUID comparison.
		if existing == id {
			// 3. Concatenate the slices before and after the removed index.
			return append(ids[:i], ids[i+1:]...)
		}
	}
	// 4. Return the original slice if the target was not found.
	return ids
}


