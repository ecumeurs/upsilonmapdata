// Package position defines the 3D coordinate system for the tactical grid.
// It provides arithmetic operations, distance calculations, and spatial utilities.
// @spec-link [[mapdata_grid_standard]]
// @spec-link [[mapdata_3d_grid]]
package position

import (
	"fmt"
	"math"

	"github.com/ecumeurs/upsilontools/tools"
)

// Position represents a discrete 3D coordinate in the Upsilon battlefield.
// X and Y represent the horizontal plane, while Z represents the vertical height.
type Position struct {
	X int
	Y int
	Z int
}

// String returns a human-readable representation of the position (X, Y, Z).
func (p Position) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.X, p.Y, p.Z)
}

// New creates a new Position instance with the specified X, Y, and Z coordinates.
// It is the standard factory function for defining discrete points in the 3D grid.
func New(x, y, z int) Position {
	// Initialize and return the static position structure.
	return Position{X: x, Y: y, Z: z}
}

// Equals returns true if both positions have identical X, Y, and Z coordinates.
func (p Position) Equals(p2 Position) bool {
	return p.X == p2.X && p.Y == p2.Y && p.Z == p2.Z
}

// SameXY returns true if both positions have identical X and Y coordinates, ignoring height.
func (p Position) SameXY(p2 Position) bool {
	return p.X == p2.X && p.Y == p2.Y
}

// Distance calculates the Manhattan distance (L1 norm) between two positions in 3D space.
func (p Position) Distance(p2 Position) int {
	return tools.Abs(p.X-p2.X) + tools.Abs(p.Y-p2.Y) + tools.Abs(p.Z-p2.Z)
}

// RandomPosition generates a random position within the specified width, length, and height bounds.
// The resulting coordinates will be in the range [0, dimension).
func RandomPosition(width, length, height int) Position {
	// 1. Generate random offsets for each axis independently.
	// 2. Return the combined coordinate as a new Position.
	return New(tools.RandomInt(0, width), tools.RandomInt(0, length), tools.RandomInt(0, height))
}

// RandomPositionInArea generates a random position within a specific sub-area of the grid.
// It uses a starting offset (x, y, z) and spans the provided dimensions.
func RandomPositionInArea(width, length, height, x, y, z int) Position {
	// 1. Calculate random offsets relative to the provided base coordinates.
	// 2. This allows for localized random placement within larger maps.
	return New(tools.RandomInt(x, x+width), tools.RandomInt(y, y+length), tools.RandomInt(z, z+height))
}


// RandomBorderPosition generates a random position located on one of the four edges of the grid.
// It is useful for spawning enemies or environmental effects at the map boundaries.
// The selection is randomized between the four cardinal edges (North, South, East, West).
func RandomBorderPosition(width, length, height int) Position {
	// 1. Roll a random integer to select one of the four border segments.
	r := tools.RandomInt(0, 4)
	switch r {
	case 0:
		// North border: Y is fixed at 0.
		return New(tools.RandomInt(0, width), 0, tools.RandomInt(0, height))
	case 1:
		// South border: Y is fixed at maximum length.
		return New(tools.RandomInt(0, width), length-1, tools.RandomInt(0, height))
	case 2:
		// West border: X is fixed at 0.
		return New(0, tools.RandomInt(0, length), tools.RandomInt(0, height))
	}
	// 3. East border: X is fixed at maximum width.
	return New(width-1, tools.RandomInt(0, length), tools.RandomInt(0, height))
}

// Add returns a new position that is the vector sum of p and p2.
func (p Position) Add(p2 Position) Position {
	return New(p.X+p2.X, p.Y+p2.Y, p.Z+p2.Z)
}

// ContainedIn returns true if the position p is present in the provided slice of positions.
// It performs a linear search using the Equals method for coordinate comparison.
func ContainedIn(p Position, positions []Position) bool {
	// 1. Iterate through every position in the target slice.
	for _, pos := range positions {
		// 2. Perform exact coordinate matching.
		if p.Equals(pos) {
			// 3. Match found: early exit with success.
			return true
		}
	}
	// 4. Exhausted all possibilities: match not found.
	return false
}


// Substract returns a new position that is the vector difference (p - p2).
func (p Position) Substract(p2 Position) Position {
	return New(p.X-p2.X, p.Y-p2.Y, p.Z-p2.Z)
}

// IsAdjacent returns true if p2 is immediately adjacent to p (orthogonally) on the X or Y axis,
// provided the height difference (Z) does not exceed the allowedJump limit.
func (p Position) IsAdjacent(p2 Position, allowedJump int) bool {
	return (p.X == p2.X && tools.Abs(p.Y-p2.Y) == 1 && tools.Abs(p.Z-p2.Z) <= allowedJump) ||
		(p.Y == p2.Y && tools.Abs(p.X-p2.X) == 1 && tools.Abs(p.Z-p2.Z) <= allowedJump)
}

// AngleTo calculates the horizontal angle (in degrees) from p to p2.
// North is 0°, East is 90°, South is 180°, and West is 270°.
func (p Position) AngleTo(p2 Position) int {
	dx := p2.X - p.X
	dy := p2.Y - p.Y
	angle := int(math.Atan2(float64(dy), float64(dx)) * 180 / 3.14)
	if angle < 0 {
		angle += 360
	}
	return angle
}


