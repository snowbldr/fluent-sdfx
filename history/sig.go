package history

//go:generate go run ./internal/opsgen/cmd

// sig is what the source knows about a function that reflection does not:
// its parameter names and the synopsis of its doc comment.
type sig struct {
	params []string
	doc    string
}
