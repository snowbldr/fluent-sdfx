// redecimate reads a binary STL (typically a raw, undecimated render) and
// runs the full DecimateMesh pipeline on it — planar snap included — writing
// the result as binary STL. This lets decimation changes be iterated on a
// saved mesh in seconds instead of re-rendering the SDF every time.
//
// Usage:
//
//	go run ./tools/redecimate <in.stl> <out.stl> [removeRatio=0.9] [maxErrMM=0.05]
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"

	flrender "github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/sdfx/render"
	"github.com/snowbldr/sdfx/sdf"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: redecimate <in.stl> <out.stl> [removeRatio=0.9] [maxErrMM=0.05]")
		os.Exit(2)
	}
	removeRatio := 0.9
	maxErr := flrender.DefaultDecimateMaxError
	if len(os.Args) >= 4 {
		removeRatio, _ = strconv.ParseFloat(os.Args[3], 64)
	}
	if len(os.Args) >= 5 {
		maxErr, _ = strconv.ParseFloat(os.Args[4], 64)
	}

	mesh, err := readBinarySTL(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "read:", err)
		os.Exit(1)
	}
	fmt.Printf("%s: %d triangles", os.Args[1], len(mesh))
	out := flrender.DecimateMesh(mesh, 1-removeRatio, maxErr)
	fmt.Println()
	if err := render.SaveSTL(os.Args[2], out); err != nil {
		fmt.Fprintln(os.Stderr, "write:", err)
		os.Exit(1)
	}
}

func readBinarySTL(path string) ([]*sdf.Triangle3, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	header := make([]byte, 84)
	if _, err := f.Read(header); err != nil {
		return nil, err
	}
	n := binary.LittleEndian.Uint32(header[80:])
	tris := make([]*sdf.Triangle3, 0, n)
	buf := make([]byte, 50)
	for i := uint32(0); i < n; i++ {
		if _, err := f.Read(buf); err != nil {
			return nil, err
		}
		tri := &sdf.Triangle3{}
		for j := 0; j < 3; j++ {
			off := 12 + j*12
			tri[j].X = float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[off:])))
			tri[j].Y = float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[off+4:])))
			tri[j].Z = float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[off+8:])))
		}
		tris = append(tris, tri)
	}
	return tris, nil
}
