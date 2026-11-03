// Package cli is the history command line, as a function so a program
// that registers its own part functions can offer the same commands:
//
//	func main() {
//		history.Register("my.Bracket", Bracket, "A bracket.", "width", "height")
//		os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
//	}
//
// Commands:
//
//	check   part.yaml               every step's status, as JSON; exit 1 if any failed or is blocked
//	build   part.yaml -o part.stl   write the bodies as one STL or 3MF (by extension)
//	export  part.yaml -o part.go    write the Go that builds the same part
//	palette                         every op, its doc and its parameters, as JSON
//	anchors part.yaml               every anchor a new step could use, as JSON
package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/snowbldr/fluent-sdfx/history"
)

// Run runs one command and returns the exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: history check|build|export|palette|anchors [file.yaml] [flags]")
		return 2
	}
	cmd, rest := args[0], args[1:]
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fs.SetOutput(stderr)
	out := fs.String("o", "", "output file")
	cells := fs.Float64("cells-per-mm", 4, "build: mesh density")
	withMain := fs.Bool("main", true, "export: add a main that writes each body to an STL")
	pkg := fs.String("package", "main", "export: package name")
	// The file comes first, flags after: history build part.yaml -o part.stl.
	var file string
	if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
		file, rest = rest[0], rest[1:]
	}
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if file == "" && fs.NArg() > 0 {
		file = fs.Arg(0)
	}
	fail := func(err error) int {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if cmd == "palette" {
		return emit(stdout, history.Palette())
	}
	if file == "" {
		fmt.Fprintf(stderr, "history %s: give the history file\n", cmd)
		return 2
	}
	h, err := history.Load(file)
	if err != nil {
		return fail(err)
	}
	switch cmd {
	case "check":
		r := history.NewRuntime(history.Options{}).Eval(h)
		if code := emit(stdout, r); code != 0 {
			return code
		}
		if r.Err() != nil {
			return 1
		}
		return 0
	case "anchors":
		return emit(stdout, history.NewRuntime(history.Options{}).Eval(h).Anchors())
	case "build":
		if *out == "" {
			return fail(fmt.Errorf("history build: give -o part.stl or -o part.3mf"))
		}
		r := history.NewRuntime(history.Options{}).Eval(h)
		if err := r.Err(); err != nil {
			return fail(err)
		}
		s := r.Solid()
		if s == nil {
			return fail(fmt.Errorf("history build: the history makes no solid"))
		}
		switch strings.ToLower(filepath.Ext(*out)) {
		case ".stl":
			s.STL(*out, *cells)
		case ".3mf":
			s.ThreeMF(*out, *cells)
		default:
			return fail(fmt.Errorf("history build: %s: write .stl or .3mf", *out))
		}
		return 0
	case "export":
		src, err := history.ExportGo(h, history.ExportOptions{Package: *pkg, Main: *withMain, CellsPerMM: *cells, Source: filepath.Base(file)})
		if err != nil {
			return fail(err)
		}
		if *out == "" {
			stdout.Write(src)
			return 0
		}
		if err := os.WriteFile(*out, src, 0o644); err != nil {
			return fail(err)
		}
		return 0
	}
	fmt.Fprintf(stderr, "history: no command %s; want check, build, export, palette or anchors\n", cmd)
	return 2
}

func emit(w io.Writer, v any) int {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
