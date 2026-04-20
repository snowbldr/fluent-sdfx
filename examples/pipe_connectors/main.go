package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
)

const name = "sch40:1"
const units = "mm"
const length = 40.0

func makePipeConnector(outputPath string, config [6]bool) {
	obj.StdPipeConnector3D(name, units, length, config).STL(outputPath, 3.0)
}

func main() {
	// 2-way
	makePipeConnector("pipe_connector_2a.stl", [6]bool{false, false, false, false, true, true})
	// 2-way
	makePipeConnector("pipe_connector_2b.stl", [6]bool{true, false, false, false, true, false})
	// 3-way
	makePipeConnector("pipe_connector_3a.stl", [6]bool{true, false, false, false, true, true})
	// 3-way
	makePipeConnector("pipe_connector_3b.stl", [6]bool{true, false, true, false, true, false})
	// 4-way
	makePipeConnector("pipe_connector_4a.stl", [6]bool{true, true, true, true, false, false})
	// 4-way
	makePipeConnector("pipe_connector_4b.stl", [6]bool{true, false, true, true, true, false})
	// 5-way
	makePipeConnector("pipe_connector_5a.stl", [6]bool{true, true, true, true, true, false})
}
