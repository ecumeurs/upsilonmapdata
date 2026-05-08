// Package pattern defines spatial shapes used for selection and effect application in the grid.
// These patterns are used by skills, area-of-effect calculations, and map generation.
// This file specifically focuses on 2D pattern operations.
// @spec-link [[mapdata_grid_standard]]
// @spec-link [[mapdata_3d_grid]]
package pattern

import (
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
)

// Pattern2D represents a collection of 2D positions (Z fixed at 0) relative to an origin.
// It is optimized for flat tactical selection on a single ground layer.
type Pattern2D []position.Position

// Square2D creates a rectangular 2D pattern with specified extents from origin.
// It generates a solid block of tiles on the flat X-Y plane (Z=0).
func Square2D(width, length int) Pattern2D {
	p := Pattern2D{}
	// 1. Construct the flat rectangle layer by layer on the X-axis.
	// 2. We iterate from negative width to positive width around the origin.
	for x := -width; x <= width; x++ {
		// 3. Delegate the longitudinal population to a helper.
		p = p.fillSquare2DY(x, length)
	}
	// 4. Return the complete 2D rectangular pattern.
	return p
}

// fillSquare2DY is a helper to populate the Y-axis for 2D squares, reducing nesting.
func (p Pattern2D) fillSquare2DY(x, length int) Pattern2D {
	// 1. Iterate through the full length of the 2D rectangle.
	// 2. This spans from negative length to positive length.
	for y := -length; y <= length; y++ {
		// 3. Add the flat position (Z=0) to the pattern slice.
		p = append(p, position.New(x, y, 0))
	}
	// 4. Stack for this X-coordinate complete.
	return p
}

// Circle2D creates a circular 2D pattern within the specified radius.
// It uses a 2D Euclidean distance check (x^2 + y^2 <= r^2).
func Circle2D(radius int) Pattern2D {
	p := Pattern2D{}
	// 1. Iterate through the bounding square of the circle on the X-axis.
	// 2. This ensures we check every possible coordinate in the circle's vicinity.
	for x := -radius; x <= radius; x++ {
		// 3. Delegate the Y-axis check to a specialized helper.
		p = p.fillCircle2DY(x, radius)
	}
	// 4. Return the fully calculated circular 2D set.
	return p
}

// fillCircle2DY is a helper to populate the Y-axis for 2D circles, reducing nesting.
func (p Pattern2D) fillCircle2DY(x, radius int) Pattern2D {
	// 1. Check every Y coordinate within the radius bounds.
	// 2. This completes the 2D plane scan for the circle.
	for y := -radius; y <= radius; y++ {
		// 3. Only include points within the circular boundary (Euclidean check).
		// 4. This ensures the resulting pattern is a smooth disk.
		if x*x+y*y <= radius*radius {
			p = append(p, position.New(x, y, 0))
		}
	}
	// 5. Circular arc segment for this X-coordinate complete.
	return p
}

// Line2D creates a linear sequence of points along the X-axis in 2D space.
// It starts at the local origin and extends forward.
func Line2D(length int) Pattern2D {
	p := Pattern2D{}
	// 1. Generate the line sequence on a single axis.
	// 2. We iterate from 0 up to the specified length.
	for i := 0; i < length; i++ {
		// 3. Add points sequentially along the positive X-axis.
		// 4. Ensure Z is always fixed at 0 for 2D compatibility.
		p = append(p, position.New(i, 0, 0))
	}
	// 5. Line generation finished successfully.
	return p
}

// Neighbours2D returns the four orthogonal adjacent tiles in a 2D grid.
// These are the North, South, East, and West neighbors.
func Neighbours2D() Pattern2D {
	// 1. Define the cardinal directions (North, South, East, West).
	// 2. These represent the immediate movement options on a flat grid.
	// 3. Each position is exactly one tile away from the origin (0,0).
	// 4. Z is fixed at 0 to ensure these are flat grid neighbors.
	return Pattern2D{
		position.New(0, 1, 0),
		position.New(0, -1, 0),
		position.New(1, 0, 0),
		position.New(-1, 0, 0),
	}
}

// PathTo2D generates a sequence of 2D positions from origin to target.
// It follows a greedy pathfinding approach to minimize distance.
func PathTo2D(pos position.Position) Pattern2D {
	p := Pattern2D{}
	curr := pos
	// 1. Trace the path from the destination back to the source.
	// 2. We use an infinite loop that breaks upon reaching origin.
	for {
		// 3. Terminate once we hit the local origin (0,0,0).
		if curr.X == 0 && curr.Y == 0 {
			break
		}
		// 4. Append the current step to our path collection.
		p = append(p, curr)
		// 5. Calculate the next greedy step closer to home.
		curr = nextPathStep2D(curr)
	}
	// 6. Reverse the sequence so it reads from source to target.
	return reversePattern2D(p)
}

// nextPathStep2D determines the next tile to move to when approaching the origin in 2D.
func nextPathStep2D(pos position.Position) position.Position {
	// 1. Greedy selection: move along the axis with the largest remaining distance.
	// 2. This ensures the path stays as direct as possible on the X/Y plane.
	// 3. Compare X and Y magnitudes to decide the next axis to increment/decrement.
	if pos.X != 0 && tools.Abs(pos.X) > tools.Abs(pos.Y) {
		// 4. Step closer on the X-axis.
		return stepX(pos)
	}
	// 5. Step closer on the Y-axis.
	return stepY(pos)
}

// reversePattern2D is a utility to reverse the order of a 2D position slice.
// It flips the backward path so it runs from the source to the target destination.
func reversePattern2D(p Pattern2D) Pattern2D {
	// 1. Swap elements from both ends toward the middle.
	// 2. Standard linear time complexity O(n).
	// 3. We use a double-iterator for the in-place reversal.
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
	// 4. Return the correctly ordered 2D sequence.
	return p
}

// Enlarge adds a neighborhood of tiles around every point in a 2D pattern.
// This is used for creating larger AOE zones from a base point or line.
func (p Pattern2D) Enlarge(radius int) Pattern2D {
	result := Pattern2D{}
	// 1. Grow each individual point by the specified radius.
	// 2. We iterate through every seed position in the existing pattern.
	for _, pos := range p {
		// 3. Merge results into the final growing shape via helper.
		result = result.enlargeAt2D(pos, radius)
	}
	// 4. Expanded 2D pattern ready for use.
	return result
}

// enlargeAt2D is a helper to grow a specific 2D point, reducing nesting.
func (p Pattern2D) enlargeAt2D(pos position.Position, radius int) Pattern2D {
	// 1. Iterate through the enlargement bounds on the X-axis.
	// 2. We scan from negative radius to positive radius around the point.
	for x := -radius; x <= radius; x++ {
		// 3. Delegate the longitudinal expansion to a helper.
		p = p.enlargeAt2DX(pos, x, radius)
	}
	// 4. Single point growth cycle complete.
	return p
}

// enlargeAt2DX is a helper to iterate the Y-axis for 2D enlargement, reducing nesting.
func (p Pattern2D) enlargeAt2DX(pos position.Position, x, radius int) Pattern2D {
	// 1. Iterate through the longitudinal expansion range.
	// 2. We scan from negative radius to positive radius around the point.
	for y := -radius; y <= radius; y++ {
		// 3. Calculate world-space target and ensure it is unique in the result.
		// 4. We use ContainedIn to avoid adding the same cell multiple times.
		target := pos.Add(position.New(x, y, 0))
		if !position.ContainedIn(target, p) {
			p = append(p, target)
		}
	}
	// 5. Finished enlargement for this specific X-offset.
	return p
}

// EnlargeVarying is a 2D variant that uses randomized radii for growth.
// This produces more natural, irregular boundaries for terrain or clouds.
func (p Pattern2D) EnlargeVarying(r tools.IntRange) Pattern2D {
	result := Pattern2D{}
	// 1. Apply growth with a unique random radius for each seed tile.
	// 2. This randomization happens independently for every point.
	for _, pos := range p {
		// 3. Use the random range to pick a radius for this specific point.
		result = result.enlargeAt2D(pos, r.Random())
	}
	// 4. Organic, jittered 2D shape finished.
	return result
}

// Apply2D translates a relative 2D pattern into absolute world-space coordinates.
// It offsets every point in the pattern by the provided start position.
func (p Pattern2D) Apply2D(startPos position.Position) []position.Position {
	result := []position.Position{}
	// 1. Translate every relative 2D point by the provided start position.
	// 2. This converts the pattern into a set of world coordinates.
	for _, pos := range p {
		// 3. Absolute coordinate calculation.
		// 4. The origin's Z value is preserved as the base height.
		result = append(result, startPos.Add(pos))
	}
	// 5. Absolute position set returned for engine consumption.
	return result
}
