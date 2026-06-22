// Package grid provides visualization utilities for the tactical grid using obelisk.js.
// It generates HTML representations that can be used for debugging or map previews.
// @spec-link [[rule_mapdata_grid_standard]]
// @spec-link [[entity_mapdata_3d_grid]]
package grid

import (
	"fmt"

	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// lighterColor calculates a color value that scales linearly with the height.
// Higher cells are rendered with lighter colors to provide visual depth in 2D views.
func lighterColor(height, maxheight, basecolor, maxcolor int) int {
	// Linear interpolation between base and max color based on height ratio.
	return basecolor + (maxcolor-basecolor)*height/maxheight
}

// hexColor formats RGB components into an "0xRRGGBB" hex string compatible with obelisk.js.
func hexColor(r, g, b int) string {
	// Use standard 6-digit hex format with 0x prefix.
	return fmt.Sprintf("0x%02x%02x%02x", r, g, b)
}

// generateCellAsObeliskCube returns a JavaScript snippet representing a 3D cube for the given cell.
// It maps the internal CellType to specific colors and applies height-based lighting.
func (g *Grid) generateCellAsObeliskCube(p position.Position) string {
	res := ""
	// 1. Define the 3D point for obelisk.js using the cell coordinates (scaled by 20).
	res += fmt.Sprintf("position = new obelisk.Point3D(%d, %d, %d);\n", p.X*20, p.Y*20, p.Z*20)
	
	// 2. Select the color based on CellType and current Z-height.
	switch g.Cells[p].Type {
	case cell.Ground:
		if g.Cells[p].IsOccupied() {
			// Occupied tiles are grayed out.
			res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 0, 80), lighterColor(p.Z, g.Height, 0, 80), lighterColor(p.Z, g.Height, 0, 80)) + ");\n"
		} else {
			// Walkable tiles are green with height-based brightness.
			res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 0, 180), 255, lighterColor(p.Z, g.Height, 0, 180)) + ");\n"
		}
	case cell.Obstacle:
		// Obstacles are reddish-brown.
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, lighterColor(p.Z, g.Height, 80, 160), lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	case cell.Water:
		// Water is blue.
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 80, 160), lighterColor(p.Z, g.Height, 50, 130), 255) + ");\n"
	case cell.Dirt:
		// Dirt represents the sub-layers.
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(150, 155, lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	case cell.Debug:
		// Debug tiles are purple.
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, lighterColor(p.Z, g.Height, 50, 130), 255) + ");\n"
	case cell.Debug2:
		// Debug2 tiles are yellow.
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, 255, lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	}
	
	// 3. Instantiate the cube and call the renderer.
	res += "cube = new obelisk.Cube(dimension, color, false);\n"
	res += "pixelView.renderObject(cube, position);\n"
	return res
}

// GenerateHTML returns an HTML representation of the grid using obelisk.js.
// This is used for debugging and visualization of generated maps.
func (g *Grid) GenerateHTML() string {
	// Header including the CDN link for obelisk.js.
	res := "<html><head><script src=\"https://unpkg.com/obelisk.js@1.2.2/build/obelisk.min.js\"></script></script></head>\n<body><canvas id=\"canvas-demo\" width=\"5000\" height=\"5000\"></canvas>\n"
	res += "<script>var canvas = document.getElementById('canvas-demo');\n"

	// Viewport setup for isometric rendering.
	res += "var point = new obelisk.Point(1000, 500);var pixelView = new obelisk.PixelView(canvas,point);\n"

	// Dimension initialization.
	res += "// create dimension instance \n"
	res += "var dimension = new obelisk.CubeDimension(20, 20, 20);\n"
	res += "// create color instance \n"
	res += "var color = new obelisk.CubeColor().getByHorizontalColor(obelisk.ColorPattern.BLUE);\n"
	res += "var cube = new obelisk.Cube(dimension, color, true);\n"
	res += "var position = new obelisk.Point3D(0, 0, 0);\n"

	// 4. Render all cubes layer by layer (Z then Y then X).
	for z := 0; z < g.Height; z++ {
		res += g.generateHTMLRow(z)
	}

	res += "</script></body></html>"
	return res
}

// generateHTMLRow is an internal helper to reduce nesting during HTML generation.
func (g *Grid) generateHTMLRow(z int) string {
	res := ""
	for y := 0; y < g.Length; y++ {
		for x := 0; x < g.Width; x++ {
			pos := position.New(x, y, z)
			// Only render cells that exist in the map.
			if _, ok := g.Cells[pos]; ok {
				res += g.generateCellAsObeliskCube(pos)
			}
		}
	}
	return res
}

