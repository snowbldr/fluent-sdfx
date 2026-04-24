package shape

import (
	flrender "github.com/snowbldr/fluent-sdfx/render"
)

// ToDXF renders the shape to a DXF file using the quadtree marching-squares renderer.
func (s *Shape) ToDXF(path string, meshCells int) {
	flrender.ToDXF(s, path, meshCells)
}

// ToSVG renders the shape to an SVG file using the quadtree marching-squares renderer.
func (s *Shape) ToSVG(path string, meshCells int) {
	flrender.ToSVG(s, path, meshCells)
}

// ToPNG rasterizes the shape to a PNG file of the given pixel dimensions,
// centered on the shape's bounding box.
func (s *Shape) ToPNG(path string, width, height int) {
	flrender.ToPNG(s, path, s.Bounds(), width, height)
}

// ToPNGBox rasterizes the shape to a PNG file for the given bounding box and pixel dimensions.
func (s *Shape) ToPNGBox(path string, bb Box2, width, height int) {
	flrender.ToPNG(s, path, bb, width, height)
}
