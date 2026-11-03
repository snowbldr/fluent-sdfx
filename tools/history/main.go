// Command history evaluates, builds and exports YAML part histories.
// See package history/cli for the commands.
//
//	go run ./tools/history check part.yaml
//	go run ./tools/history build part.yaml -o part.stl
package main

import (
	"os"

	"github.com/snowbldr/fluent-sdfx/history/cli"
)

func main() { os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr)) }
