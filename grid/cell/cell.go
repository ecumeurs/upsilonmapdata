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

// Cell is a cell in the grid
type Cell struct {
	Type     CellType
	EntityID uuid.UUID
	Position position.Position
}

// NewCell
func NewCell(t CellType, p position.Position) *Cell {
	return &Cell{
		Type:     t,
		Position: p,
		EntityID: uuid.Nil,
	}
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
