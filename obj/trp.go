package obj

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/obj"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// TruncRectPyramidParms configures a truncated rectangular pyramid.
type TruncRectPyramidParms struct {
	Size         v3.Vec  // size of truncated pyramid
	BaseAngleDeg float64 // base angle of pyramid (degrees, like the rest of the API)
	BaseRadius   float64 // base corner radius
	RoundRadius  float64 // edge rounding radius
}

func (p *TruncRectPyramidParms) toSDF() *obj.TruncRectPyramidParms {
	return &obj.TruncRectPyramidParms{
		Size:        v3sdf.Vec(p.Size),
		BaseAngle:   p.BaseAngleDeg * math.Pi / 180,
		BaseRadius:  p.BaseRadius,
		RoundRadius: p.RoundRadius,
	}
}

// TruncRectPyramid3D returns a truncated rectangular pyramid.
func TruncRectPyramid3D(p TruncRectPyramidParms) *solid.Solid {
	s, err := obj.TruncRectPyramid3D(p.toSDF())
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
