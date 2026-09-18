// Command score compares one candidate package against a task's reference.
//
//	go run ./bench/cmd/score -task bench/tasks/pill-boss-length -cand bench/tasks/pill-boss-length/asfailed
//
// The candidate directory must be a Go package inside this module that
// exports Build() *solid.Solid. Prints the score.Result as JSON.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/snowbldr/fluent-sdfx/bench/score"
)

func main() {
	task := flag.String("task", "", "task directory")
	cand := flag.String("cand", "", "candidate package directory (default: <task>/reference)")
	work := flag.String("work", "", "scratch dir for the generated scorer (default: bench/_runs/_score/<task>)")
	root := flag.String("root", ".", "repo root")
	flag.Parse()
	if *task == "" {
		flag.Usage()
		os.Exit(2)
	}
	if *cand == "" {
		*cand = filepath.Join(*task, "reference")
	}
	if *work == "" {
		*work = filepath.Join(*root, "bench", "_runs", "_score", filepath.Base(*task))
	}
	r, err := score.Score(*root, *task, *cand, *work)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	b, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(b))
	if !r.Passed {
		os.Exit(1)
	}
}
