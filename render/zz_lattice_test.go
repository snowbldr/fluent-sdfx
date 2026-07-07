package render

import (
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/snowbldr/fluent-sdfx/render/meshopt"
)

// Direct reproduction: run the real snap on the raw mold mesh, crop the +X
// wall face (which refuses to collapse in production), and simplify just
// that region standalone. Set FLSDFX_RAW_STL to the raw mesh path.
func TestMeshoptRealFaceCrop(t *testing.T) {
	path := os.Getenv("FLSDFX_RAW_STL")
	if path == "" {
		t.Skip("set FLSDFX_RAW_STL")
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	header := make([]byte, 84)
	f.Read(header)
	n := int(binary.LittleEndian.Uint32(header[80:]))
	xmin := float32(53.5)
	if v := os.Getenv("FLSDFX_CROP_XMIN"); v != "" {
		f, _ := strconv.ParseFloat(v, 32)
		xmin = float32(f)
	}
	buf := make([]byte, 50)
	verts := make([]float32, 0, n*9)
	for i := 0; i < n; i++ {
		f.Read(buf)
		keep := false
		var tri [9]float32
		for j := 0; j < 3; j++ {
			off := 12 + j*12
			for k := 0; k < 3; k++ {
				tri[j*3+k] = math.Float32frombits(binary.LittleEndian.Uint32(buf[off+k*4:]))
			}
			if tri[j*3] > xmin { // crop threshold
				keep = true
			}
		}
		if keep {
			verts = append(verts, tri[:]...)
		}
	}
	count := len(verts) / 9
	t.Logf("cropped %d tris near +X wall", count)

	// crop extent for relative error scaling
	minv, maxv := [3]float32{99e9, 99e9, 99e9}, [3]float32{-99e9, -99e9, -99e9}
	for i := 0; i < count*3; i++ {
		for k := 0; k < 3; k++ {
			v := verts[i*3+k]
			if v < minv[k] {
				minv[k] = v
			}
			if v > maxv[k] {
				maxv[k] = v
			}
		}
	}
	extent := 0.0
	for k := 0; k < 3; k++ {
		if e := float64(maxv[k] - minv[k]); e > extent {
			extent = e
		}
	}
	t.Logf("crop extent %.1f", extent)

	// raw, unsnapped — at the crop's own relative scale and at the FULL
	// mesh's relative scale (same absolute budget, different normalization)
	_, c1, e1 := meshopt.Simplify(append([]float32(nil), verts...), count, count/100, float32(0.01/extent))
	t.Logf("raw @rel(crop)=%.3g:  %d → %d (err %.3g)", 0.01/extent, count, c1, e1)

	// same content, 4x the size: four translated copies in one soup. If
	// meshopt degrades on this, the stall is a pure scale effect.
	if os.Getenv("FLSDFX_TILE") != "" {
		tiled := make([]float32, 0, len(verts)*4)
		for k := 0; k < 4; k++ {
			dy := float32(k * 150)
			for i := 0; i < len(verts); i += 3 {
				tiled = append(tiled, verts[i], verts[i+1]+dy, verts[i+2])
			}
		}
		tcount := len(tiled) / 9
		// same ABSOLUTE budget: tiled extent grows to ~3*150+ext in y
		text := extent + 450
		if 450+extent < 450 {
			text = 450
		}
		_, c4, e4 := meshopt.Simplify(tiled, tcount, tcount/100, float32(0.01/text))
		t.Logf("tiled 4x @0.01mm abs: %d → %d (%.1f%% kept, err %.3g)", tcount, c4, 100*float64(c4)/float64(tcount), e4)
	}

}
