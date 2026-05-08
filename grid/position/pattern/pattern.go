// Package pattern defines spatial shapes used for selection and effect application in the grid.
// These patterns are used by skills, area-of-effect calculations, and map generation.
// @spec-link [[mapdata_grid_standard]]
// @spec-link [[mapdata_3d_grid]]
package pattern

import (
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
)

// Pattern represents a collection of 3D positions relative to an origin.
// It is the primary structure for defining area-of-effect and multi-tile selections.
type Pattern []position.Position

// Contains returns true if the pattern, when centered at origin, includes the target position.
// This is done by translating each pattern point by the origin and comparing it to pos.
func (p Pattern) Contains(origin, pos position.Position) bool {
	// 1. Begin iteration through all relative points in the pattern.
	for _, pt := range p {
		// 2. Calculate the absolute position in world-space.
		// 3. Check for exact coordinate equality against the target.
		if origin.Add(pt).Equals(pos) {
			return true
		}
	}
	// 4. No match found after exhaustive search.
	return false
}

// ContainsAny returns true if the pattern includes at least one of the provided positions.
// This is useful for checking if a player is standing in any part of a multi-tile zone.
func (p Pattern) ContainsAny(origin position.Position, positions []position.Position) bool {
	// 1. Loop through each position provided in the input slice.
	for _, pos := range positions {
		// 2. Delegate to the standard Contains method for individual checks.
		if p.Contains(origin, pos) {
			return true
		}
	}
	// 3. None of the input positions are within the pattern bounds.
	return false
}

// ContainsAll returns true if the pattern includes all of the provided positions.
// Use this to verify if an entire unit or structure is fully covered by an effect.
func (p Pattern) ContainsAll(origin position.Position, positions []position.Position) bool {
	// 1. Every single position must be validated against the pattern.
	for _, pos := range positions {
		// 2. If even one position fails, the entire check fails.
		if !p.Contains(origin, pos) {
			return false
		}
	}
	// 3. Success: All positions are verified as contained.
	return true
}

// Single returns a 3D pattern containing only the origin point (0, 0, 0).
// This is used for single-target skills that still require a pattern interface.
func Single() Pattern {
	// Return a slice with a single position at origin.
	return Pattern{position.New(0, 0, 0)}
}

// Circle creates a spherical 3D pattern within the specified integer radius.
// It uses a simple Euclidean distance check (x^2 + y^2 + z^2 <= r^2).
func Circle(radius int) Pattern {
	p := Pattern{}
	// 1. Iterate through the bounding cube of the sphere on the X-axis.
	for x := -radius; x <= radius; x++ {
		// 2. Delegate the Y and Z axis population to helper methods.
		p = p.fillCircleY(x, radius)
	}
	// 3. Return the fully populated spherical pattern.
	return p
}

// fillCircleY is a helper to populate the Y-axis of a 3D circle, reducing nesting.
func (p Pattern) fillCircleY(x, radius int) Pattern {
	// 1. Iterate through the Y-range for the current X coordinate.
	for y := -radius; y <= radius; y++ {
		// 2. Further delegate to the Z-axis population helper.
		p = p.fillCircleZ(x, y, radius)
	}
	// 3. Return the updated pattern slice.
	return p
}

// fillCircleZ is a helper to populate the Z-axis of a 3D circle, reducing nesting.
func (p Pattern) fillCircleZ(x, y, radius int) Pattern {
	// 1. Iterate through the Z-range to complete the cubic volume check.
	for z := -radius; z <= radius; z++ {
		// 2. Include points where the distance to origin is within the radius.
		// 3. Uses the standard Euclidean distance formula: x^2 + y^2 + z^2.
		if x*x+y*y+z*z <= radius*radius {
			p = append(p, position.New(x, y, z))
		}
	}
	// 4. Final slice for this Z-stack is returned.
	return p
}

// Square creates a cubic 3D pattern with the specified dimensions from origin.
// width, length, and height are the half-extents of the cube.
func Square(width, length, height int) Pattern {
	p := Pattern{}
	// 1. Construct the cube layer by layer starting from the X-axis.
	for x := -width; x <= width; x++ {
		// 2. Use helper methods to manage vertical and longitudinal stacks.
		p = p.fillSquareY(x, length, height)
	}
	// 3. Return the complete 3D rectangular volume.
	return p
}

// fillSquareY is a helper to populate the Y-axis of a 3D square, reducing nesting.
func (p Pattern) fillSquareY(x, length, height int) Pattern {
	// 1. Iterate through the full length specified for the pattern.
	for y := -length; y <= length; y++ {
		// 2. Delegate to the Z-axis populator for the vertical dimension.
		p = p.fillSquareZ(x, y, height)
	}
	// 3. Return the updated pattern for this (x, y) column.
	return p
}

// fillSquareZ is a helper to populate the Z-axis of a 3D square, reducing nesting.
func (p Pattern) fillSquareZ(x, y, height int) Pattern {
	// 1. Iterate through the full vertical height specified.
	for z := -height; z <= height; z++ {
		// 2. All points within the rectangular bounds are unconditionally added.
		p = append(p, position.New(x, y, z))
	}
	// 3. Column population complete.
	return p
}

// Line creates a 1D line pattern starting at (0,0,0) and extending along the positive X-axis.
func Line(length int) Pattern {
	p := Pattern{}
	// 1. Generate a linear sequence of points along the primary axis.
	for i := 0; i < length; i++ {
		// 2. Each point is incremented by exactly one tile on the X-axis.
		p = append(p, position.New(i, 0, 0))
	}
	// 3. Return the 1D pattern sequence.
	return p
}

// Neighbours returns a 3x3x3 cubic pattern representing all cells adjacent to the origin.
func Neighbours() Pattern {
	// Adjacent cells are defined as those within a radius of 1 in all directions.
	return Square(1, 1, 1)
}

// AdjacentDirections returns the six primary cardinal 3D directions.
// These represent North, South, East, West, Up, and Down in the coordinate system.
func AdjacentDirections() Pattern {
	// 1. Explicitly list the six unit vectors for 3D cardinal movement.
	// 2. These are used for calculating adjacency and immediate pathing options.
	// 3. The vectors cover all six faces of a cubic grid cell.
	return Pattern{
		position.New(0, 0, 1),
		position.New(0, 0, -1),
		position.New(0, 1, 0),
		position.New(0, -1, 0),
		position.New(1, 0, 0),
		position.New(-1, 0, 0),
	}
}


// Apply transforms the relative pattern points into absolute world-space coordinates.
// Each point in the pattern is offset by the provided origin position.
func (p Pattern) Apply(origin position.Position) []position.Position {
	result := []position.Position{}
	// 1. Loop through the relative points defined in the pattern.
	for _, pt := range p {
		// 2. Translate every point by adding the origin coordinates.
		result = append(result, origin.Add(pt))
	}
	// 3. Return the absolute world-space positions.
	return result
}

// ApplyInArea is a bounded version of Apply that only returns positions within the grid.
// It filters out any points that would fall outside the specified width, length, or height.
func (p Pattern) ApplyInArea(origin position.Position, width, length, height int) []position.Position {
	result := []position.Position{}
	// 1. Transform each point and then immediately validate its coordinates.
	for _, pt := range p {
		pt = origin.Add(pt)
		// 2. Boundary Check: Ensure coordinates are >= 0 and < grid limits.
		if pt.X >= 0 && pt.X < width && pt.Y >= 0 && pt.Y < length && pt.Z >= 0 && pt.Z < height {
			result = append(result, pt)
		}
	}
	// 3. Return only the valid, in-bounds positions.
	return result
}

// PathTo generates a direct path of points from the origin (0,0,0) to the target.
// It greedily steps along the axis with the greatest remaining distance.
func PathTo(pos position.Position) Pattern {
	p := Pattern{}
	curr := pos
	// 1. Move backward from the target destination toward the origin point.
	for {
		// 2. Terminate the loop once we reach the local origin.
		if curr.X == 0 && curr.Y == 0 && curr.Z == 0 {
			break
		}
		// 3. Record the current step and calculate the next one.
		p = append(p, curr)
		curr = nextPathStep(curr)
	}
	// 4. The path was built backward, so we must reverse it for the final result.
	return reversePattern(p)
}

// nextPathStep calculates the next logical coordinate when moving toward origin.
func nextPathStep(pos position.Position) position.Position {
	// 1. Compare the absolute distances across all three axes.
	// 2. Prioritize moving along the axis with the largest delta to stay direct.
	if pos.X != 0 && tools.Abs(pos.X) > tools.Abs(pos.Y) {
		return stepX(pos)
	} else if pos.Y != 0 && tools.Abs(pos.Y) < tools.Abs(pos.Z) {
		return stepY(pos)
	}
	// 3. Fallback to stepping along the Z-axis.
	return stepZ(pos)
}

// stepX increments or decrements X to move it closer to zero.
func stepX(pos position.Position) position.Position {
	// 1. Check current sign and move toward origin.
	// 2. This is a unit step to maintain contiguous paths.
	if pos.X > 0 {
		pos.X--
	} else {
		pos.X++
	}
	// 3. Return the updated position.
	return pos
}

// stepY increments or decrements Y to move it closer to zero.
func stepY(pos position.Position) position.Position {
	// 1. Check current sign and move toward origin.
	// 2. This is a unit step to maintain contiguous paths.
	if pos.Y > 0 {
		pos.Y--
	} else {
		pos.Y++
	}
	// 3. Return the updated position.
	return pos
}

// stepZ increments or decrements Z to move it closer to zero.
func stepZ(pos position.Position) position.Position {
	// 1. Check current sign and move toward origin.
	// 2. This is a unit step to maintain contiguous paths.
	if pos.Z > 0 {
		pos.Z--
	} else {
		pos.Z++
	}
	// 3. Return the updated position.
	return pos
}

// reversePattern flips the order of points in a slice, making a backward path forward.
func reversePattern(p Pattern) Pattern {
	// 1. Use a standard two-pointer swap to reverse the slice in-place.
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
	// 2. Return the correctly ordered path.
	return p
}

// Enlarge creates a thicker pattern by adding a neighborhood around every existing point.
// This is used for creating larger AOE zones from a base point or line.
func (p Pattern) Enlarge(radius int) Pattern {
	result := Pattern{}
	// 1. For every point in the original pattern, grow its boundary.
	for _, pt := range p {
		// 2. Merge the resulting expansion into the final result set.
		result = result.enlargeAt(pt, radius)
	}
	// 3. Return the expanded volume.
	return result
}

// enlargeAt is a helper to add points around a specific position, reducing nesting.
func (p Pattern) enlargeAt(pos position.Position, radius int) Pattern {
	// 1. Iterate through the expansion volume on the X-axis.
	for x := -radius; x <= radius; x++ {
		// 2. Delegate the remaining dimensions to secondary helpers.
		p = p.enlargeAtX(pos, x, radius)
	}
	// 3. Return the updated aggregate pattern.
	return p
}

// enlargeAtX is a helper to iterate the Y-axis for enlargement, reducing nesting.
func (p Pattern) enlargeAtX(pos position.Position, x, radius int) Pattern {
	// 1. Iterate through the Y-range for expansion.
	for y := -radius; y <= radius; y++ {
		// 2. Delegate to the final Z-axis populator.
		p = p.enlargeAtY(pos, x, y, radius)
	}
	// 3. Return the updated pattern slice.
	return p
}

// enlargeAtY is a helper to iterate the Z-axis for enlargement, reducing nesting.
func (p Pattern) enlargeAtY(pos position.Position, x, y, radius int) Pattern {
	// 1. Iterate through the vertical Z-range for expansion.
	for z := -radius; z <= radius; z++ {
		// 2. Calculate the absolute coordinate of the new neighbor.
		target := pos.Add(position.New(x, y, z))
		// 3. Ensure uniqueness by checking for existing duplicates in the slice.
		if !position.ContainedIn(target, p) {
			p = append(p, target)
		}
	}
	// 4. Completion of vertical stack enlargement.
	return p
}

// EnlargeVarying applies Enlarge but uses a random radius for every single point.
// This produces more organic, irregular shapes for natural terrain generation.
func (p Pattern) EnlargeVarying(r tools.IntRange) Pattern {
	result := Pattern{}
	// 1. Each point in the pattern can grow by a different, randomized amount.
	for _, pt := range p {
		// 2. Use a random radius selection for each seed point.
		result = result.enlargeAt(pt, r.Random())
	}
	// 3. Return the organically grown shape.
	return result
}
