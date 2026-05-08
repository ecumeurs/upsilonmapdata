package grid

import (
	"fmt"

	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// lighterColor calculates a color value that scales with the height.
// Higher cells are rendered with lighter colors to provide visual depth in 2D views.
// @lint-ignore-documentation
func lighterColor(height, maxheight, basecolor, maxcolor int) int {
	return basecolor + (maxcolor-basecolor)*height/maxheight
}

// hexColor formats RGB components into an "0xRRGGBB" hex string.
// @lint-ignore-documentation
func hexColor(r, g, b int) string {
	return fmt.Sprintf("0x%02x%02x%02x", r, g, b)
}

// generateCellAsObeliskCube returns the cell as an obelisk cube string.
// This is used for generating the isometric 3D view of the grid.
func (g *Grid) generateCellAsObeliskCube(p position.Position) string {
	res := ""
	res += fmt.Sprintf("position = new obelisk.Point3D(%d, %d, %d);\n", p.X*20, p.Y*20, p.Z*20)
	switch g.Cells[p].Type {
	case cell.Ground:
		if g.Cells[p].IsOccupied() {
			res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 0, 80), lighterColor(p.Z, g.Height, 0, 80), lighterColor(p.Z, g.Height, 0, 80)) + ");\n"
		} else {
			res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 0, 180), 255, lighterColor(p.Z, g.Height, 0, 180)) + ");\n"
		}
	case cell.Obstacle:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, lighterColor(p.Z, g.Height, 80, 160), lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	case cell.Water:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 80, 160), lighterColor(p.Z, g.Height, 50, 130), 255) + ");\n"
	case cell.Dirt:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(150, 155, lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	case cell.Debug:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, lighterColor(p.Z, g.Height, 50, 130), 255) + ");\n"
	case cell.Debug2:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, 255, lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	}
	res += "cube = new obelisk.Cube(dimension, color, false);\n"
	res += "pixelView.renderObject(cube, position);\n"
	return res
}

// GenerateHTML returns an HTML representation of the grid using obelisk.js.
// This is used for debugging and visualization of generated maps.
func (g *Grid) GenerateHTML() string {
	res := "<html><head><script src=\"https://unpkg.com/obelisk.js@1.2.2/build/obelisk.min.js\"></script></script></head>\n<body><canvas id=\"canvas-demo\" width=\"5000\" height=\"5000\"></canvas>\n"
	res += "<script>var canvas = document.getElementById('canvas-demo');\n"

	res += "var point = new obelisk.Point(1000, 500);var pixelView = new obelisk.PixelView(canvas,point);\n"

	res += "// create dimension instance \n"
	res += "var dimension = new obelisk.CubeDimension(20, 20, 20);\n"
	res += "// create color instance \n"
	res += "var color = new obelisk.CubeColor().getByHorizontalColor(obelisk.ColorPattern.BLUE);\n"
	res += "var cube = new obelisk.Cube(dimension, color, true);\n"
	res += "var position = new obelisk.Point3D(0, 0, 0);\n"

	// render all cubes
	for z := 0; z < g.Height; z++ {
		for y := 0; y < g.Length; y++ {
			for x := 0; x < g.Width; x++ {
				pos := position.New(x, y, z)
				_, ok := g.Cells[pos]
				if !ok {
					continue
				}
				res += g.generateCellAsObeliskCube(pos)
			}
		}
	}

	res += "</script></body></html>"
	return res
}
