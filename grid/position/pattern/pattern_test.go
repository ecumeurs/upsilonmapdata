// Package pattern provides tests for spatial shape generation and manipulation.
// It ensures that geometric algorithms correctly map to grid coordinates.
// @test-link [[rule_mapdata_grid_standard]]
// @test-link [[entity_grid]]
package pattern

import (
	"testing"

	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// TestPathTo verifies that the PathTo2D algorithm produces a contiguous chain of adjacent cells
// and correctly reaches the intended destination coordinates.
// It validates the core pathfinding logic used for 2D movement and targeting.
func TestPathTo(t *testing.T) {
	// 1. Define start and end positions for the test path.
	// We use coordinates far enough apart to test the multi-step greedy logic.
	startPos := position.New(30, 20, 0)
	endPos := position.New(40, 10, 0)
	
	// 2. Calculate the relative path from origin (0,0) to the target offset.
	// The algorithm works with relative offsets to be origin-agnostic.
	path := PathTo2D(endPos.Substract(startPos))

	// 3. Adjacency Check Phase:
	// Iterate through the generated path and verify adjacency between sequential points.
	// This prevents "teleportation" or diagonal skips that violate engine rules.
	for i := 1; i < len(path); i++ {
		// Each step must be within a distance of 1 tile (plus height jump allowance).
		if !path[i].IsAdjacent(path[i-1], 2) {
			t.Errorf("Path step %d is not adjacent to the previous one", i)
		}
	}

	// 4. Logging & Visualization Phase:
	// Output the generated path for visual verification in test logs.
	t.Log("Start position:", startPos)
	// Transform relative path steps into absolute world positions for logging.
	expandedPath := path.Apply2D(startPos)
	for i := 0; i < len(expandedPath); i++ {
		// Log every single step coordinate for detailed audit trails.
		t.Logf("Step %d: %v", i, expandedPath[i])
	}
	t.Log("End position:", endPos)

	// 5. Destination Verification Phase:
	// Verify that the final point in the path exactly matches the target destination.
	// Failure here indicates a bug in the greedy step calculation or accumulation.
	if !expandedPath[len(expandedPath)-1].Equals(endPos) {
		t.Errorf("Last path step is not the end position")
	}
	
	// 6. Test successfully completed all validation phases.
}



