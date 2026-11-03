// Command cmd rewrites history/ops_gen.go from the library source.
//
//	go generate ./history/     # runs this from the history directory
package main

import (
	"flag"
	"log"
	"os"

	"github.com/snowbldr/fluent-sdfx/history/internal/opsgen"
)

func main() {
	root := flag.String("root", "..", "module root")
	out := flag.String("o", "ops_gen.go", "output file")
	flag.Parse()
	src, err := opsgen.Generate(*root)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, src, 0o644); err != nil {
		log.Fatal(err)
	}
}
