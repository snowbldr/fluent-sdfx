// Package history runs a part written as data: a YAML (or JSON) list of
// steps, each one a fluent-sdfx function or method, evaluated by
// reflection. Every click in a UI and every tool call from an agent is
// one step; deleting a step and evaluating again makes it never have
// happened; the same file exports to a Go program that builds the same
// part.
//
//	h, err := history.Load("figure.yaml")
//	rt := history.NewRuntime(history.Options{})
//	res := rt.Eval(h)          // every step's status; never fails as a whole
//	part := res.Solid()        // the bodies, as one solid
//	src, err := history.ExportGo(h, history.ExportOptions{Main: true})
//
// A step names its op (solid.Sphere, sculpt.Pull, or a method such as
// Translate on its input), its args by parameter name or by position, and
// optionally a name other steps refer to it by. Nothing else: no
// arithmetic, no control flow. Anything procedural is a Go function
// registered with Register, which a history then calls like any other op.
//
// The runtime caches each step by a hash of its definition and of every
// step it depends on, so evaluating again after an edit recomputes only
// the edited step and what follows. Provenance (by, from) is not part of
// the hash. Errors are per step: ok, error (with a message that names the
// fix), blocked (something it needs failed or does not exist), or off.
//
// Palette describes every op for a UI or an agent's tool list, from the
// same tables the evaluator uses, so the description cannot drift from
// what runs. The function table is generated from the library source
// (go generate ./history/) and a test fails when it is stale.
//
// Status: experimental. The schema is in README.md; names are
// placeholders until the owner settles them.
package history
