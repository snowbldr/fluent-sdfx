// Command selfcheck validates every task in bench/tasks:
//
//   - reference vs reference must pass (IoU 1, all probes pass);
//   - reference vs asfailed (the construction the model actually produced
//     in the source incident) must FAIL, proving the metric discriminates.
//
//	go run ./bench/cmd/selfcheck            # all tasks
//	go run ./bench/cmd/selfcheck pill-boss  # tasks whose name contains the arg
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/snowbldr/fluent-sdfx/bench/score"
)

func main() {
	root := flag.String("root", ".", "repo root")
	verbose := flag.Bool("v", false, "print full results")
	flag.Parse()
	filter := flag.Arg(0)

	tasksDir := filepath.Join(*root, "bench", "tasks")
	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), "_") && (filter == "" || strings.Contains(e.Name(), filter)) {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	fails := 0
	fmt.Printf("%-28s %-10s %7s %7s %8s %s\n", "task", "candidate", "iou", "maxdev", "probes", "verdict")
	for _, n := range names {
		task := filepath.Join(tasksDir, n)
		for _, cand := range []string{"reference", "asfailed"} {
			cdir := filepath.Join(task, cand)
			if _, err := os.Stat(cdir); err != nil {
				if cand == "asfailed" {
					continue
				}
				fmt.Printf("%-28s %-10s missing\n", n, cand)
				fails++
				continue
			}
			r, err := score.Score(*root, task, cdir, filepath.Join(*root, "bench", "_runs", "_selfcheck", n, cand))
			if err != nil {
				fmt.Printf("%-28s %-10s ERROR %v\n", n, cand, err)
				fails++
				continue
			}
			wantPass := cand == "reference"
			ok := r.Passed == wantPass && r.Built
			verdict := "ok"
			if !ok {
				verdict = "UNEXPECTED"
				fails++
			}
			if !r.Built {
				fmt.Printf("%-28s %-10s build failed: %s\n", n, cand, firstLine(r.Error))
				continue
			}
			fmt.Printf("%-28s %-10s %7.3f %7.2f %4d/%-3d %s %s\n", n, cand, r.IoU, r.MaxDev, r.ProbesPass, r.ProbesTotal, verdict, strings.Join(r.Reasons, "; "))
			if *verbose {
				for _, p := range r.Probes {
					mark := "  "
					if !p.Passed {
						mark = "!!"
					}
					fmt.Printf("    %s %-32s want=%-5v got=%-5v sdf=%.3f\n", mark, p.Name, p.Want, p.Got, p.SDF)
				}
			}
		}
	}
	if fails > 0 {
		fmt.Printf("\n%d unexpected result(s)\n", fails)
		os.Exit(1)
	}
	fmt.Printf("\nall %d task(s) ok\n", len(names))
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
