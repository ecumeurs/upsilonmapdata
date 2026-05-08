package grid

import (
	"fmt"

	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// MoveEntity moves an entity from one position to another within the grid.
// It removes the entity from the source cell and adds it to the destination cell.
// @spec-link [[mechanic_multi_entity_cell_system]]
func (g *Grid) MoveEntity(from, to position.Position, entityID uuid.UUID) error {
	// Verify that the destination position exists within the grid boundaries.
	if !g.PositionIsInGrid(to) {
		return fmt.Errorf("to position %v is not in the grid", to)
	}
	
	// Remove the entity from its current location if the cell exists.
	if c, ok := g.CellAt(from); ok {
		c.RemoveEntity(entityID)
	}
	
	// Retrieve the destination cell and add the entity to it.
	c, ok := g.CellAt(to)
	if !ok {
		return fmt.Errorf("to position %v is not in the grid", to)
	}
	c.AddEntity(entityID)
	return nil
}

// AddEntity adds an entity to a cell at the given position.
// This is typically used during entity spawning or initialization.
// @spec-link [[mechanic_multi_entity_cell_system]]
func (g *Grid) AddEntity(p position.Position, entityID uuid.UUID) {
	if c, ok := g.CellAt(p); ok {
		c.AddEntity(entityID)
	}
}

// RemoveEntity removes a specific entity from the cell at the given position.
// This is used when an entity leaves the grid or is destroyed.
// @spec-link [[mechanic_multi_entity_cell_system]]
func (g *Grid) RemoveEntity(p position.Position, entityID uuid.UUID) {
	if c, ok := g.CellAt(p); ok {
		c.RemoveEntity(entityID)
	}
}

// GetEntitiesAt returns all entity IDs currently present at the specified position.
// It returns a copy of the slice to prevent external modification of the cell state.
// @spec-link [[mechanic_multi_entity_cell_system]]
func (g *Grid) GetEntitiesAt(p position.Position) []uuid.UUID {
	if c, ok := g.CellAt(p); ok {
		result := make([]uuid.UUID, len(c.EntityIDs))
		copy(result, c.EntityIDs)
		return result
	}
	return nil
}

// IsOccupiedByOther returns true if the cell contains any entity other than the given selfID.
// This is used for collision detection and pathfinding validation.
func (g *Grid) IsOccupiedByOther(p position.Position, selfID uuid.UUID) bool {
	if c, ok := g.CellAt(p); ok {
		for _, id := range c.EntityIDs {
			if id != selfID {
				return true
			}
		}
	}
	return false
}
