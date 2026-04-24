package cell

import (
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// CellType describe all the possible types of cell
type CellType int

const (
	// WallCellType is the type of a wall cell
	Obstacle CellType = 0
	// GroundCellType is the type of a ground cell
	Ground CellType = 1
	Water  CellType = 2
	Dirt   CellType = 3
	Debug  CellType = 4
	Debug2 CellType = 5
)

// Cell is a cell in the grid.
// A cell can contain multiple entities simultaneously (characters + traps + area effects etc).
//
// @spec-link [[mechanic_multi_entity_cell_system]]
// @spec-link [[mechanic_cell_attached_effects]]
type Cell struct {
	Type     CellType
	Position position.Position

	// EntityIDs holds all entities currently in this cell.
	// Multiple entities are allowed: one character + any number of WalkThrough entities.
	EntityIDs []uuid.UUID

	// EffectIDs holds positional effects attached to this cell (traps, zones, terrain modifiers).
	// Actual effect data is stored in GameState.Effects.
	EffectIDs []uuid.UUID
}

// NewCell
func NewCell(t CellType, p position.Position) *Cell {
	return &Cell{
		Type:      t,
		Position:  p,
		EntityIDs: []uuid.UUID{},
		EffectIDs: []uuid.UUID{},
	}
}

// AddEntity adds an entity ID to this cell.
func (c *Cell) AddEntity(id uuid.UUID) {
	c.EntityIDs = append(c.EntityIDs, id)
}

// RemoveEntity removes an entity ID from this cell.
func (c *Cell) RemoveEntity(id uuid.UUID) {
	c.EntityIDs = removeID(c.EntityIDs, id)
}

// HasEntity returns true if the given entity ID is in this cell.
func (c *Cell) HasEntity(id uuid.UUID) bool {
	for _, eid := range c.EntityIDs {
		if eid == id {
			return true
		}
	}
	return false
}

// IsOccupied returns true if any entity is in this cell.
func (c *Cell) IsOccupied() bool {
	return len(c.EntityIDs) > 0
}

// AddEffect adds a positional effect ID to this cell.
func (c *Cell) AddEffect(id uuid.UUID) {
	c.EffectIDs = append(c.EffectIDs, id)
}

// RemoveEffect removes a positional effect ID from this cell.
func (c *Cell) RemoveEffect(id uuid.UUID) {
	c.EffectIDs = removeID(c.EffectIDs, id)
}

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

// removeID removes an ID from a slice, preserving order. Returns the same slice if not found.
func removeID(ids []uuid.UUID, id uuid.UUID) []uuid.UUID {
	for i, existing := range ids {
		if existing == id {
			return append(ids[:i], ids[i+1:]...)
		}
	}
	return ids
}
