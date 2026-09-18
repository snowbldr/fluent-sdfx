// Command inspect renders an STL from the standard inspection views so
// that whoever built it, person or model, can look at it before anyone
// prints it.
//
//	go run ./tools/inspect part.stl                 # six PNGs next to the STL
//	go run ./tools/inspect -out shots/ -w 900 part.stl other.stl
//	go run ./tools/inspect -views top,front part.stl
//
// Views: iso-top-front, iso-top-back, iso-bottom-front, iso-bottom-back,
// top, bottom, front, back, left, right. The default set is the four isos
// plus top and bottom, which is what caught the floating bodies, the sealed
// window and the mirrored stamp in the record this tool comes from. Z is up.
//
// Needs f3d on PATH (https://f3d.app). For a section view, cut the solid
// first (Solid.CutPlane, or validate.FitCheck's cutaway) and inspect that
// STL; f3d has no clipping in this workflow.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var views = map[string]string{
	"iso-top-front":    "1,-1,-0.8",
	"iso-top-back":     "-1,1,-0.8",
	"iso-bottom-front": "1,-1,0.8",
	"iso-bottom-back":  "-1,1,0.8",
	"top":              "0,0,-1",
	"bottom":           "0,0,1",
	"front":            "0,1,0",
	"back":             "0,-1,0",
	"left":             "1,0,0",
	"right":            "-1,0,0",
}

var defaultViews = []string{"iso-top-front", "iso-top-back", "iso-bottom-front", "iso-bottom-back", "top", "bottom"}

func main() {
	out := flag.String("out", "", "directory for PNGs (default: next to each STL)")
	width := flag.Int("w", 800, "image width in pixels")
	height := flag.Int("h", 700, "image height in pixels")
	which := flag.String("views", strings.Join(defaultViews, ","), "comma-separated view names")
	color := flag.String("color", "#8fbc8f", "surface colour")
	edges := flag.Bool("edges", false, "draw mesh edges")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if _, err := exec.LookPath("f3d"); err != nil {
		fmt.Fprintln(os.Stderr, "inspect: f3d not found on PATH; install from https://f3d.app")
		os.Exit(1)
	}
	var names []string
	for _, v := range strings.Split(*which, ",") {
		v = strings.TrimSpace(v)
		if _, ok := views[v]; !ok {
			fmt.Fprintf(os.Stderr, "inspect: unknown view %q; known: %s\n", v, strings.Join(sortedViews(), ", "))
			os.Exit(2)
		}
		names = append(names, v)
	}
	failed := 0
	for _, stl := range flag.Args() {
		dir := *out
		if dir == "" {
			dir = filepath.Dir(stl)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		base := strings.TrimSuffix(filepath.Base(stl), filepath.Ext(stl))
		for _, v := range names {
			png := filepath.Join(dir, base+"-"+v+".png")
			// Pin the look so a user's f3d config (grid, opacity, filename
			// overlay) cannot change what an inspection shows.
			args := []string{stl, "--output=" + png, fmt.Sprintf("--resolution=%d,%d", *width, *height),
				"--up=+Z", "--camera-direction=" + views[v], "--color=" + *color,
				"--grid=false", "--opacity=1", "--filename=false", "--axis=true"}
			if *edges {
				args = append(args, "--edges")
			}
			cmd := exec.Command("f3d", args...)
			if b, err := cmd.CombinedOutput(); err != nil {
				failed++
				fmt.Fprintf(os.Stderr, "inspect: %s %s: %v\n%s\n", stl, v, err, strings.TrimSpace(string(b)))
				continue
			}
			fmt.Println(png)
		}
	}
	if failed > 0 {
		os.Exit(1)
	}
}

func sortedViews() []string {
	out := make([]string, 0, len(views))
	for k := range views {
		out = append(out, k)
	}
	for i := range out {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
