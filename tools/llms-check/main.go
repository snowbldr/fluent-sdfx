// Command llms-check keeps docs/llms.txt, the agent-facing reference, in
// step with the exported API.
//
// The reference has two files. docs/llms.txt is written by hand: the
// mental model, the conventions, the traps, the recipes; it is what a model
// loads first and it must stay short. docs/llms-api.txt is generated from
// the source: one entry per exported symbol with its signature and the
// synopsis of its doc comment, so a model that needs the whole surface can
// load it and always finds every symbol that exists, described in the words
// the code itself carries.
//
// Checks, all of which fail the build:
//
//   - docs/llms-api.txt differs from a fresh generation (run -gen);
//   - an exported symbol has no doc comment, so its index entry would be
//     empty (write one; the doc comment is the source of truth);
//   - a code block in docs/llms.txt references pkg.Name( or .Method( that
//     does not exist (the narrative has outlived the code).
//
// Symbols that are deliberately undocumented go in docs/llms-ignore.txt,
// one per line as pkg.Name or pkg.Type.Method, a reason after '#'.
//
//	go run ./tools/llms-check            # check; exit 1 on drift
//	go run ./tools/llms-check -gen       # rewrite docs/llms-api.txt
//	go run ./tools/llms-check -undoc     # list exported symbols lacking doc comments
//
// Run by `make build`, and therefore by the pre-commit hook.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// packages whose exported surface the doc must cover, keyed by the
// identifier an importer would use, in the order they appear in the index.
var packages = []struct{ name, dir, blurb string }{
	{"solid", "solid", "3D solids: primitives, transforms, booleans, blends, patterns, modifiers, anchors, output."},
	{"shape", "shape", "2D shapes: primitives, builders, transforms, booleans, threads, cams, text, and the 2D-to-3D operations."},
	{"layout", "layout", "Position lists for Multi: polar, grid, line, corners."},
	{"plane", "plane", "Planes for cross-sections."},
	{"obj", "obj", "Parametric parts: bolts, nuts, panels, gears, keyed holes, pipes, and more. Each takes a parameter struct."},
	{"validate", "validate", "Geometry checks for tests: probes, mesh stats, printability."},
	{"render", "render", "Output formats and renderer constructors."},
	{"mesh", "mesh", "Triangle-mesh utilities."},
	{"units", "units", "Constants and unit conversion."},
	{"v3", "vec/v3", "3D vectors and boxes."},
	{"v2", "vec/v2", "2D vectors and boxes."},
	{"p2", "vec/p2", "Polar 2D points."},
	{"v3i", "vec/v3i", "Integer 3D vectors."},
	{"v2i", "vec/v2i", "Integer 2D vectors."},
}

type symbol struct {
	pkg, recv, name, kind string // kind: func, method, type, const, var
	sig, doc              string
}

func (s symbol) key() string {
	if s.recv != "" {
		return s.pkg + "." + s.recv + "." + s.name
	}
	return s.pkg + "." + s.name
}

func main() {
	root := flag.String("root", ".", "repo root")
	doc := flag.String("doc", "docs/llms.txt", "hand-written agent-facing reference (narrative)")
	api := flag.String("api", "docs/llms-api.txt", "generated API index")
	ignore := flag.String("ignore", "docs/llms-ignore.txt", "deliberately undocumented symbols")
	gen := flag.Bool("gen", false, "rewrite the generated API index file")
	undoc := flag.Bool("undoc", false, "list exported symbols that have no doc comment, and exit")
	flag.Parse()

	syms, err := exported(*root)
	must(err)
	ignored := readIgnore(filepath.Join(*root, *ignore))
	var live []symbol
	for _, s := range syms {
		if ignored[s.key()] || ignored[s.pkg+"."+s.name] || ignored[s.pkg+".*"] {
			continue
		}
		live = append(live, s)
	}

	if *undoc {
		n := 0
		for _, s := range live {
			if s.doc == "" {
				fmt.Printf("%-8s %s\n", s.kind, s.key())
				n++
			}
		}
		fmt.Printf("%d undocumented of %d exported\n", n, len(live))
		return
	}

	docPath := filepath.Join(*root, *doc)
	text, err := os.ReadFile(docPath)
	must(err)
	narrative := string(text)
	index := generate(live)
	apiPath := filepath.Join(*root, *api)

	if *gen {
		must(os.WriteFile(apiPath, []byte(index), 0o644))
		fmt.Printf("llms-check: wrote %s (%d symbols)\n", *api, len(live))
		return
	}

	fail := false

	// 1. Doc comments. The index is only as good as they are.
	var missingDoc []symbol
	for _, s := range live {
		if s.doc == "" {
			missingDoc = append(missingDoc, s)
		}
	}
	if len(missingDoc) > 0 {
		fail = true
		fmt.Printf("llms-check: %d exported symbol(s) have no doc comment:\n", len(missingDoc))
		for _, s := range missingDoc {
			fmt.Printf("  %-8s %s\n", s.kind, s.key())
		}
		fmt.Printf("  (write one; or list the symbol in %s with a reason)\n", *ignore)
	}

	// 2. Index freshness.
	if current, err := os.ReadFile(apiPath); err != nil || string(current) != index {
		fail = true
		fmt.Printf("llms-check: %s is stale or missing; run: go run ./tools/llms-check -gen\n", *api)
	}

	// 3. Stale references in the narrative's code blocks.
	known := map[string]bool{}
	methods := map[string]bool{}
	pkgSet := map[string]bool{}
	for _, s := range syms {
		pkgSet[s.pkg] = true
		if s.kind == "method" {
			methods[s.name] = true
		} else {
			known[s.pkg+"."+s.name] = true
		}
	}
	var names []string
	for p := range pkgSet {
		names = append(names, p)
	}
	sort.Strings(names)
	pkgAlt := strings.Join(names, "|")
	callRe := regexp.MustCompile(`\b(` + pkgAlt + `)\.([A-Z][A-Za-z0-9_]*)\b`)
	methRe := regexp.MustCompile(`\.([A-Z][A-Za-z0-9_]*)\(`)
	var stale []string
	seen := map[string]bool{}
	for _, block := range codeBlocks(narrative) {
		for _, m := range callRe.FindAllStringSubmatch(block, -1) {
			key := m[1] + "." + m[2]
			if !known[key] && !seen[key] {
				seen[key] = true
				stale = append(stale, key)
			}
		}
		for _, m := range methRe.FindAllStringSubmatch(block, -1) {
			name := m[1]
			k := "." + name + "("
			if methods[name] || seen[k] {
				continue
			}
			// pkg.Func( is covered above; only flag genuine method-looking calls,
			// and not calls into the standard library.
			if regexp.MustCompile(`\b(` + pkgAlt + `|fmt|os|math|strings|testing|errors|time|sort|log)\.` + regexp.QuoteMeta(name) + `\(`).MatchString(block) {
				continue
			}
			seen[k] = true
			stale = append(stale, k)
		}
	}
	if len(stale) > 0 {
		fail = true
		sort.Strings(stale)
		fmt.Printf("llms-check: %d identifier(s) in %s code blocks do not exist:\n", len(stale), *doc)
		for _, s := range stale {
			fmt.Printf("  %s\n", s)
		}
	}

	if fail {
		os.Exit(1)
	}
	fmt.Printf("llms-check: %s references nothing stale; %s indexes all %d exported symbols\n", *doc, *api, len(live))
}

// generate renders the API index: one section per package, one line per
// symbol, methods grouped under their receiver type.
func generate(syms []symbol) string {
	var b strings.Builder
	b.WriteString("# fluent-sdfx API index\n\n")
	b.WriteString("Every exported symbol, generated from the source by tools/llms-check; do not edit by hand. Read docs/llms.txt first for the mental model and the traps; use this file to look up a signature or find a symbol you did not know existed. Each entry is the signature followed by the synopsis of the symbol's doc comment. Import paths are `github.com/snowbldr/fluent-sdfx/<package>`; the vector packages are `vec/v3`, `vec/v2`, `vec/p2`, `vec/v3i`, `vec/v2i`.\n")
	for _, p := range packages {
		var funcs, types, consts []symbol
		byRecv := map[string][]symbol{}
		for _, s := range syms {
			if s.pkg != p.name {
				continue
			}
			switch s.kind {
			case "method":
				byRecv[s.recv] = append(byRecv[s.recv], s)
			case "type":
				types = append(types, s)
			case "func":
				funcs = append(funcs, s)
			default:
				consts = append(consts, s)
			}
		}
		if len(funcs)+len(types)+len(consts)+len(byRecv) == 0 {
			continue
		}
		fmt.Fprintf(&b, "\n### %s — %s\n", p.name, p.blurb)
		if len(consts) > 0 {
			b.WriteString("\nConstants and variables:\n")
			for _, s := range consts {
				fmt.Fprintf(&b, "- `%s.%s` — %s\n", s.pkg, s.name, s.doc)
			}
		}
		if len(funcs) > 0 {
			b.WriteString("\nFunctions:\n")
			for _, s := range funcs {
				fmt.Fprintf(&b, "- `%s.%s%s` — %s\n", s.pkg, s.name, s.sig, s.doc)
			}
		}
		for _, t := range types {
			fmt.Fprintf(&b, "\n`%s.%s` — %s\n", t.pkg, t.name, t.doc)
			ms := byRecv[t.name]
			delete(byRecv, t.name)
			for _, m := range ms {
				fmt.Fprintf(&b, "- `.%s%s` — %s\n", m.name, m.sig, m.doc)
			}
		}
		// Methods whose receiver type is declared elsewhere (aliases etc).
		var rest []string
		for r := range byRecv {
			rest = append(rest, r)
		}
		sort.Strings(rest)
		for _, r := range rest {
			fmt.Fprintf(&b, "\nMethods on `%s.%s`:\n", p.name, r)
			for _, m := range byRecv[r] {
				fmt.Fprintf(&b, "- `.%s%s` — %s\n", m.name, m.sig, m.doc)
			}
		}
	}
	return b.String()
}

func codeBlocks(body string) []string {
	re := regexp.MustCompile("(?s)```[a-z]*\\n(.*?)```")
	var out []string
	for _, m := range re.FindAllStringSubmatch(body, -1) {
		out = append(out, m[1])
	}
	return out
}

func readIgnore(path string) map[string]bool {
	out := map[string]bool{}
	f, err := os.Open(path)
	if err != nil {
		return out
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		if line != "" {
			out[line] = true
		}
	}
	return out
}

// exported parses every public package and returns its exported symbols
// with signatures and first-sentence docs, sorted by key.
func exported(root string) ([]symbol, error) {
	var out []symbol
	for _, p := range packages {
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, filepath.Join(root, p.dir), func(fi os.FileInfo) bool {
			return !strings.HasSuffix(fi.Name(), "_test.go")
		}, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p.dir, err)
		}
		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				for _, d := range file.Decls {
					switch d := d.(type) {
					case *ast.FuncDecl:
						if !d.Name.IsExported() {
							continue
						}
						s := symbol{pkg: p.name, name: d.Name.Name, kind: "func", sig: funcSig(fset, d.Type), doc: firstSentence(d.Doc)}
						if d.Recv != nil && len(d.Recv.List) > 0 {
							recv := recvName(d.Recv.List[0].Type)
							if recv == "" || !ast.IsExported(recv) {
								continue
							}
							s.recv, s.kind = recv, "method"
						}
						out = append(out, s)
					case *ast.GenDecl:
						for _, spec := range d.Specs {
							switch sp := spec.(type) {
							case *ast.TypeSpec:
								if !sp.Name.IsExported() {
									continue
								}
								doc := firstSentence(sp.Doc)
								if doc == "" {
									doc = firstSentence(d.Doc)
								}
								out = append(out, symbol{pkg: p.name, name: sp.Name.Name, kind: "type", doc: doc})
							case *ast.ValueSpec:
								kind := "var"
								if d.Tok == token.CONST {
									kind = "const"
								}
								doc := firstSentence(sp.Doc)
								if doc == "" {
									doc = firstSentence(d.Doc)
								}
								for _, n := range sp.Names {
									if n.IsExported() {
										out = append(out, symbol{pkg: p.name, name: n.Name, kind: kind, doc: doc})
									}
								}
							}
						}
					}
				}
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].key() < out[j].key() })
	return out, nil
}

// funcSig renders "(params) results" for a FuncType.
func funcSig(fset *token.FileSet, ft *ast.FuncType) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, ft)
	s := buf.String()
	s = strings.TrimPrefix(s, "func")
	return strings.Join(strings.Fields(s), " ")
}

// firstSentence returns the synopsis of a doc comment (go/doc's rule: the
// first sentence, where "e.g." and similar do not end one), or "".
func firstSentence(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}
	// go/doc ends a sentence at ". " unless the period follows a single
	// capital letter, so "e.g. " and "i.e. " would truncate it. Shield them.
	t := cg.Text()
	for _, ab := range []string{"e.g.", "i.e.", "etc.", "vs.", "approx.", "min.", "max."} {
		t = strings.ReplaceAll(t, ab+" ", ab+"\x00")
	}
	return strings.Join(strings.Fields(strings.ReplaceAll(doc.Synopsis(t), "\x00", " ")), " ")
}

func recvName(t ast.Expr) string {
	switch t := t.(type) {
	case *ast.StarExpr:
		return recvName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return recvName(t.X)
	}
	return ""
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
