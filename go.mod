module github.com/snowbldr/fluent-sdfx

go 1.25.0

require github.com/deadsy/sdfx v0.0.1

// The snowbldr fork of sdfx carries upstream fixes (Screw3D accuracy,
// parallel octree renderer, bbox-pruned unions, flat-polygon SDF)
// that fluent-sdfx depends on. Imports stay as github.com/deadsy/sdfx
// so source stays portable if the fork is merged upstream.
//
// TODO: once the sdfx fork is pushed and tagged v0.0.1, swap this to
//   replace github.com/deadsy/sdfx => github.com/snowbldr/sdfx v0.0.1
replace github.com/deadsy/sdfx => /Users/r/sdfx

require (
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/dhconnelly/rtreego v1.2.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/hpinc/go3mf v0.24.2 // indirect
	github.com/llgcode/draw2d v0.0.0-20260213073409-1c39bbefe083 // indirect
	github.com/qmuntal/opc v0.7.12 // indirect
	github.com/yofu/dxf v0.0.0-20250806094206-f3988c7f0176 // indirect
	golang.org/x/image v0.38.0 // indirect
)
