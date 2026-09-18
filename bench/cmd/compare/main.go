// Command compare prints a task-by-condition matrix across several runs.
//
//	go run ./bench/cmd/compare -cols "opus/none=baseline-opus-none+fix-opus-none,opus/full=baseline-opus-full+fix-opus-full"
//
// Each column is a label and a '+'-joined list of run directories (relative
// to bench/_runs or absolute). Later runs in a list override earlier ones for
// the same task, which is how a re-run of one task after a harness fix
// replaces the original result. Cells show pass, FAIL, or the stage reached,
// plus the number of model replies and whether a question was asked.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type result struct {
	Task            string  `json:"task"`
	Category        string  `json:"category"`
	ExpectsQuestion bool    `json:"expects_question"`
	AskedFirst      bool    `json:"asked_first"`
	Questions       int     `json:"questions"`
	Turns           int     `json:"turns"`
	CostUSD         float64 `json:"cost_usd"`
	Error           string  `json:"error"`
	Score           struct {
		Compiled bool    `json:"compiled"`
		Built    bool    `json:"built"`
		Passed   bool    `json:"passed"`
		IoU      float64 `json:"iou"`
	} `json:"score"`
}

func main() {
	cols := flag.String("cols", "", "label=run+run,label=run,...")
	root := flag.String("root", ".", "repo root")
	flag.Parse()
	if *cols == "" {
		flag.Usage()
		os.Exit(2)
	}
	runsDir := filepath.Join(*root, "bench", "_runs")

	type column struct {
		label string
		res   map[string]result
	}
	var columns []column
	tasks := map[string]string{} // task -> category
	for _, spec := range strings.Split(*cols, ",") {
		label, runs, ok := strings.Cut(strings.TrimSpace(spec), "=")
		if !ok {
			fmt.Fprintf(os.Stderr, "bad column spec %q\n", spec)
			os.Exit(2)
		}
		c := column{label: label, res: map[string]result{}}
		for _, run := range strings.Split(runs, "+") {
			dir := run
			if !filepath.IsAbs(dir) {
				dir = filepath.Join(runsDir, run)
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", dir, err)
				os.Exit(1)
			}
			for _, e := range entries {
				b, err := os.ReadFile(filepath.Join(dir, e.Name(), "result.json"))
				if err != nil {
					continue
				}
				var r result
				if json.Unmarshal(b, &r) != nil {
					continue
				}
				c.res[r.Task] = r
				tasks[r.Task] = r.Category
			}
		}
		columns = append(columns, c)
	}

	names := make([]string, 0, len(tasks))
	for t := range tasks {
		names = append(names, t)
	}
	sort.Strings(names)

	fmt.Printf("| task | category |")
	for _, c := range columns {
		fmt.Printf(" %s |", c.label)
	}
	fmt.Printf("\n|---|---|")
	for range columns {
		fmt.Printf("---|")
	}
	fmt.Println()
	for _, t := range names {
		fmt.Printf("| %s | %s |", t, tasks[t])
		for _, c := range columns {
			r, ok := c.res[t]
			if !ok {
				fmt.Printf(" - |")
				continue
			}
			fmt.Printf(" %s |", cell(r))
		}
		fmt.Println()
	}
	fmt.Printf("| **passed** | |")
	for _, c := range columns {
		p, n, cost, turns, ask, askN := 0, 0, 0.0, 0, 0, 0
		for _, r := range c.res {
			n++
			if r.Score.Passed {
				p++
			}
			cost += r.CostUSD
			turns += r.Turns
			if r.ExpectsQuestion {
				askN++
				if r.AskedFirst {
					ask++
				}
			}
		}
		fmt.Printf(" **%d/%d** asked %d/%d, %.2f turns, $%.2f |", p, n, ask, askN, float64(turns)/float64(max(n, 1)), cost)
	}
	fmt.Println()
}

func cell(r result) string {
	s := "pass"
	switch {
	case r.Error != "" && !r.Score.Built:
		s = "error"
	case !r.Score.Compiled:
		s = "compile"
	case !r.Score.Built:
		s = "build"
	case !r.Score.Passed:
		s = "FAIL"
	}
	extra := ""
	if r.Questions > 0 {
		extra = " ?"
	}
	if r.Turns > 1 {
		extra += fmt.Sprintf(" t%d", r.Turns)
	}
	return s + extra
}
