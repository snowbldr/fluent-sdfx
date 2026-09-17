// Command run drives a model over benchmark tasks with `claude -p`, plays a
// scripted user when the model asks a clarifying question, scores each
// result, and writes a summary.
//
//	go run ./bench/cmd/run -model opus -tools none
//	go run ./bench/cmd/run -model fable -tools full -tasks pill,teardrop -jobs 2
//
// Two tool conditions:
//
//	none  the model gets no tools and must reply with the finished file in a
//	      ```go fenced block. Measures pure translation.
//	full  the model works in a workspace inside this module with Read, Write,
//	      Edit, Glob, Grep and Bash limited to go/f3d/ls/cat, so it can
//	      compile, render and probe its own work. Measures translation plus
//	      self-verification.
//
// Clarify protocol: the model may reply with a line starting QUESTION:. The
// runner answers from the task's answers.md via a small "user" model (or
// "Your call, use a sensible default." when there is no hidden spec) and
// resumes the session. Turns are counted; whether the model asked before
// building is recorded and compared against the task's expects_question.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/snowbldr/fluent-sdfx/bench/score"
)

type claudeOut struct {
	Result       string  `json:"result"`
	SessionID    string  `json:"session_id"`
	IsError      bool    `json:"is_error"`
	NumTurns     int     `json:"num_turns"`
	TotalCostUSD float64 `json:"total_cost_usd"`
	DurationMS   int     `json:"duration_ms"`
	StopReason   string  `json:"stop_reason"`
}

type turn struct {
	Role string `json:"role"` // user | model | scripted-user
	Text string `json:"text"`
}

type taskRun struct {
	Task            string       `json:"task"`
	Category        string       `json:"category"`
	Mode            string       `json:"mode"`
	ExpectsQuestion bool         `json:"expects_question"`
	AskedFirst      bool         `json:"asked_first"`
	Questions       int          `json:"questions"`
	Turns           int          `json:"turns"` // model replies consumed
	CostUSD         float64      `json:"cost_usd"`
	ElapsedS        float64      `json:"elapsed_s"`
	Transcript      []turn       `json:"transcript"`
	Score           score.Result `json:"score"`
	Error           string       `json:"error,omitempty"`
}

var (
	model     = flag.String("model", "opus", "model for the candidate (claude -p --model)")
	userModel = flag.String("user-model", "sonnet", "model that plays the scripted user")
	tools     = flag.String("tools", "none", "tool condition: none | full")
	tasksArg  = flag.String("tasks", "", "comma-separated substrings of task names to run (default all)")
	maxTurns  = flag.Int("turns", 6, "max model replies per task")
	jobs      = flag.Int("jobs", 2, "tasks in parallel")
	refDoc    = flag.String("ref", "docs/llms.txt", "reference doc appended to the system prompt (empty to disable)")
	outDir    = flag.String("out", "", "run directory (default bench/_runs/<time>-<model>-<tools>)")
	root      = flag.String("root", ".", "repo root")
	callTO    = flag.Duration("timeout", 15*time.Minute, "timeout per claude call")
	dry       = flag.Bool("dry", false, "prepare workspaces and prompts, do not call the model")
	rescore   = flag.String("rescore", "", "re-score an existing run directory with the current scorer (no model calls) and rewrite its summary")
)

var fence = regexp.MustCompile("(?s)```(?:go|golang)?\\s*\\n(.*?)```")

func main() {
	flag.Parse()
	rootAbs, _ := filepath.Abs(*root)
	if *rescore != "" {
		rescoreRun(rootAbs, *rescore)
		return
	}
	if *outDir == "" {
		*outDir = filepath.Join(rootAbs, "bench", "_runs", fmt.Sprintf("%s-%s-%s", time.Now().Format("20060102-150405"), *model, *tools))
	}
	outAbs, _ := filepath.Abs(*outDir)
	must(os.MkdirAll(outAbs, 0o755))

	refText := ""
	if *refDoc != "" {
		b, err := os.ReadFile(filepath.Join(rootAbs, *refDoc))
		must(err)
		refText = string(b)
	}

	names := selectTasks(rootAbs, *tasksArg)
	if len(names) == 0 {
		fmt.Fprintln(os.Stderr, "no tasks matched")
		os.Exit(1)
	}
	fmt.Printf("run: %s\nmodel=%s tools=%s tasks=%d jobs=%d\n\n", outAbs, *model, *tools, len(names), *jobs)

	results := make([]taskRun, len(names))
	var wg sync.WaitGroup
	sem := make(chan struct{}, *jobs)
	for i, n := range names {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, n string) {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = runTask(rootAbs, outAbs, n, refText)
			r := results[i]
			fmt.Printf("%-28s %s\n", n, oneLine(r))
		}(i, n)
	}
	wg.Wait()

	writeSummary(outAbs, results)
}

func selectTasks(rootAbs, filter string) []string {
	entries, err := os.ReadDir(filepath.Join(rootAbs, "bench", "tasks"))
	must(err)
	var subs []string
	if filter != "" {
		subs = strings.Split(filter, ",")
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), "_") {
			continue
		}
		if len(subs) == 0 {
			names = append(names, e.Name())
			continue
		}
		for _, s := range subs {
			if strings.Contains(e.Name(), strings.TrimSpace(s)) {
				names = append(names, e.Name())
				break
			}
		}
	}
	sort.Strings(names)
	return names
}

func runTask(rootAbs, outAbs, name, refText string) taskRun {
	start := time.Now()
	taskDir := filepath.Join(rootAbs, "bench", "tasks", name)
	spec, err := score.LoadSpec(taskDir)
	if err != nil {
		return taskRun{Task: name, Error: err.Error()}
	}
	tr := taskRun{Task: name, Category: spec.Category, Mode: spec.Mode, ExpectsQuestion: spec.ExpectsQuestion}

	ws := filepath.Join(outAbs, name)
	candDir := filepath.Join(ws, "candidate")
	must(os.MkdirAll(candDir, 0o755))

	// Seed the candidate: start code for modify tasks, a stub otherwise.
	startCode := ""
	if b, err := os.ReadFile(filepath.Join(taskDir, "start", "start.go")); err == nil {
		startCode = strings.Replace(string(b), "package start", "package candidate", 1)
	} else {
		startCode = "package candidate\n\nimport \"github.com/snowbldr/fluent-sdfx/solid\"\n\n// Build returns the part.\nfunc Build() *solid.Solid {\n\tpanic(\"not implemented\")\n}\n"
	}
	must(os.WriteFile(filepath.Join(candDir, "candidate.go"), []byte(startCode), 0o644))

	promptB, err := os.ReadFile(filepath.Join(taskDir, "prompt.md"))
	must(err)
	prompt := strings.TrimSpace(string(promptB))
	answers := ""
	if b, err := os.ReadFile(filepath.Join(taskDir, "answers.md")); err == nil {
		answers = string(b)
	}

	sys := systemPrompt(*tools, refText)
	must(os.WriteFile(filepath.Join(ws, "system.md"), []byte(sys), 0o644))

	userMsg := prompt
	if spec.Mode == "modify" {
		userMsg += "\n\nHere is the current code, candidate/candidate.go (package candidate). Modify it:\n\n```go\n" + startCode + "```\n"
	} else {
		userMsg += "\n\nStart from scratch; the file candidate/candidate.go currently holds a stub.\n"
	}
	tr.Transcript = append(tr.Transcript, turn{"user", userMsg})
	must(os.WriteFile(filepath.Join(ws, "prompt.txt"), []byte(userMsg), 0o644))
	if *dry {
		return tr
	}

	sessionID := ""
	var final string
	for tr.Turns < *maxTurns {
		out, raw, err := callClaude(ws, sys, userMsg, sessionID, *model, *tools)
		tr.Turns++
		if err != nil {
			tr.Error = fmt.Sprintf("claude call failed: %v\n%s", err, raw)
			break
		}
		sessionID = out.SessionID
		tr.CostUSD += out.TotalCostUSD
		tr.Transcript = append(tr.Transcript, turn{"model", out.Result})
		if q, ok := extractQuestion(out.Result); ok {
			tr.Questions++
			if tr.Turns == 1 {
				tr.AskedFirst = true
			}
			ans := scriptedUser(ws, prompt, answers, q)
			tr.Transcript = append(tr.Transcript, turn{"scripted-user", ans})
			userMsg = ans
			continue
		}
		final = out.Result
		break
	}
	if tr.Error == "" && final == "" && tr.Turns >= *maxTurns {
		tr.Error = "turn budget exhausted without a final answer"
	}

	// Materialise the candidate.
	if *tools == "none" && final != "" {
		if m := fence.FindAllStringSubmatch(final, -1); len(m) > 0 {
			code := m[len(m)-1][1]
			code = regexp.MustCompile(`(?m)^package \w+`).ReplaceAllString(code, "package candidate")
			must(os.WriteFile(filepath.Join(candDir, "candidate.go"), []byte(code), 0o644))
		} else if tr.Error == "" {
			tr.Error = "no fenced go block in final reply"
		}
	}

	saveJSON(filepath.Join(ws, "transcript.json"), tr.Transcript)
	res, err := score.Score(rootAbs, taskDir, candDir, filepath.Join(ws, "_scorer"))
	if err != nil {
		tr.Error = strings.TrimSpace(tr.Error + "\nscore: " + err.Error())
	} else {
		tr.Score = res
	}
	tr.ElapsedS = time.Since(start).Seconds()
	saveJSON(filepath.Join(ws, "result.json"), tr)
	return tr
}

func systemPrompt(tools, refText string) string {
	var b strings.Builder
	b.WriteString(`You are translating a person's description of a mechanical part into fluent-sdfx Go code.

Deliverable: the file candidate/candidate.go, package candidate, exporting func Build() *solid.Solid. Import paths are github.com/snowbldr/fluent-sdfx/solid, /shape, /layout, /obj, /plane, /vec/v3, /vec/v2. Units are millimetres, angles are degrees, Z is up.

Protocol:
- If the description leaves out something you need and cannot reasonably default (which face, which orientation, a dimension or a component detail that changes the geometry), ask before building. Reply with ONLY a line starting with QUESTION: followed by your question(s), nothing else. The person will answer and you continue. Do not ask about things a sensible default settles, and do not ask more than necessary.
- Do exactly what was asked. Do not add features, tapers, fillets, clearances or improvements that were not requested, and change nothing the description did not ask you to change.
`)
	if tools == "none" {
		b.WriteString("- You have no tools. When you are done, reply with the complete file in a single ```go fenced block. Put nothing after the block.\n")
	} else {
		b.WriteString("- Your working directory is a Go package directory inside the fluent-sdfx module. Write the file to candidate/candidate.go, check it compiles with `go vet ./candidate`, and verify the geometry however you like (write a throwaway main under ./check/, evaluate the SDF at probe points, render an STL and inspect it with f3d --output). When you are satisfied, reply with the single word DONE.\n")
	}
	if refText != "" {
		b.WriteString("\n--- fluent-sdfx reference ---\n")
		b.WriteString(refText)
	}
	return b.String()
}

func callClaude(ws, sys, msg, sessionID, model, tools string) (claudeOut, string, error) {
	args := []string{"-p", msg, "--model", model, "--output-format", "json"}
	if sessionID == "" {
		sysPath := filepath.Join(ws, "system.md")
		args = append(args, "--append-system-prompt-file", sysPath)
	} else {
		args = append(args, "--resume", sessionID)
	}
	if tools == "none" {
		args = append(args, "--tools", "")
	} else {
		args = append(args, "--allowedTools", "Read,Write,Edit,Glob,Grep,Bash(go:*),Bash(f3d:*),Bash(ls:*),Bash(cat:*),Bash(mkdir:*)",
			"--permission-mode", "acceptEdits")
	}
	ctx, cancel := context.WithTimeout(context.Background(), *callTO)
	defer cancel()
	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Dir = ws
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	raw := stdout.String() + stderr.String()
	var out claudeOut
	if jerr := json.Unmarshal(stdout.Bytes(), &out); jerr != nil {
		if err == nil {
			err = jerr
		}
		return out, raw, err
	}
	if out.IsError {
		return out, raw, fmt.Errorf("claude returned is_error: %s", out.Result)
	}
	return out, raw, nil
}

func extractQuestion(reply string) (string, bool) {
	t := strings.TrimSpace(reply)
	if strings.HasPrefix(t, "QUESTION:") {
		return strings.TrimSpace(strings.TrimPrefix(t, "QUESTION:")), true
	}
	// Tolerate a short preamble before the marker, but not a reply that
	// also contains code.
	if i := strings.Index(t, "QUESTION:"); i >= 0 && i < 200 && !strings.Contains(t, "```") {
		return strings.TrimSpace(t[i+len("QUESTION:"):]), true
	}
	return "", false
}

func scriptedUser(ws, prompt, answers, question string) string {
	if strings.TrimSpace(answers) == "" {
		return "Your call, use a sensible default."
	}
	sys := "You are the person who wrote the request below. Someone building the part for you has asked one or more questions. The hidden spec below describes the part you want; it is the source of truth. Answer every sub-question the spec covers, plainly and with the spec's numbers, even if the question is framed around another part or drawing you mentioned (a socket, a mockup, an existing nub): those were only comparisons, and the spec's numbers are what you want. Only for a sub-question the spec genuinely does not cover, say: Your call, use a sensible default. Do not volunteer information that was not asked about. Do not write code.\n\n--- your request ---\n" + prompt + "\n\n--- hidden spec ---\n" + answers
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "claude", "-p", question, "--model", *userModel, "--tools", "", "--output-format", "json", "--no-session-persistence", "--system-prompt", sys)
	cmd.Dir = ws
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "Your call, use a sensible default."
	}
	var out claudeOut
	if json.Unmarshal(stdout.Bytes(), &out) != nil || out.Result == "" {
		return "Your call, use a sensible default."
	}
	return strings.TrimSpace(out.Result)
}

func oneLine(r taskRun) string {
	if r.Error != "" && !r.Score.Built {
		return "ERROR " + firstLine(r.Error)
	}
	s := r.Score
	stage := "pass"
	if !s.Compiled {
		stage = "compile-fail"
	} else if !s.Built {
		stage = "build-fail"
	} else if !s.Passed {
		stage = "FAIL"
	}
	ask := ""
	if r.ExpectsQuestion {
		ask = fmt.Sprintf(" asked=%v(expected)", r.AskedFirst)
	} else if r.Questions > 0 {
		ask = fmt.Sprintf(" asked=%d(unexpected)", r.Questions)
	}
	return fmt.Sprintf("%-12s iou=%.3f dev=%.2f probes=%d/%d turns=%d%s $%.2f %s",
		stage, s.IoU, s.MaxDev, s.ProbesPass, s.ProbesTotal, r.Turns, ask, r.CostUSD, strings.Join(s.Reasons, "; "))
}

func writeSummary(outAbs string, rs []taskRun) {
	saveJSON(filepath.Join(outAbs, "summary.json"), rs)
	var b strings.Builder
	fmt.Fprintf(&b, "# bench run\n\nmodel=%s tools=%s tasks=%d\n\n", *model, *tools, len(rs))
	fmt.Fprintf(&b, "| task | category | stage | iou | max_dev | probes | turns | asked | expected | cost |\n|---|---|---|---|---|---|---|---|---|---|\n")
	pass, compiled, built := 0, 0, 0
	var turns, cost float64
	clarifyN, clarifyHit, falseAsk := 0, 0, 0
	for _, r := range rs {
		s := r.Score
		stage := "pass"
		switch {
		case r.Error != "" && !s.Built:
			stage = "error"
		case !s.Compiled:
			stage = "compile"
		case !s.Built:
			stage = "build"
		case !s.Passed:
			stage = "fail"
		}
		if s.Compiled {
			compiled++
		}
		if s.Built {
			built++
		}
		if s.Passed {
			pass++
		}
		turns += float64(r.Turns)
		cost += r.CostUSD
		if r.ExpectsQuestion {
			clarifyN++
			if r.AskedFirst {
				clarifyHit++
			}
		} else if r.Questions > 0 {
			falseAsk++
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %.3f | %.2f | %d/%d | %d | %d | %v | $%.2f |\n",
			r.Task, r.Category, stage, s.IoU, s.MaxDev, s.ProbesPass, s.ProbesTotal, r.Turns, r.Questions, r.ExpectsQuestion, r.CostUSD)
	}
	n := float64(len(rs))
	fmt.Fprintf(&b, "\n- passed: %d/%d\n- compiled: %d, built: %d\n- mean turns: %.2f\n- clarify tasks asked first: %d/%d; unexpected questions on other tasks: %d\n- total cost: $%.2f\n",
		pass, len(rs), compiled, built, turns/n, clarifyHit, clarifyN, falseAsk, cost)
	must(os.WriteFile(filepath.Join(outAbs, "summary.md"), []byte(b.String()), 0o644))
	fmt.Printf("\npassed %d/%d, mean turns %.2f, clarify %d/%d, cost $%.2f\nsummary: %s\n",
		pass, len(rs), turns/n, clarifyHit, clarifyN, cost, filepath.Join(outAbs, "summary.md"))
}

func saveJSON(path string, v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	must(os.WriteFile(path, b, 0o644))
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rescoreRun re-runs the scorer over every task workspace in an existing run
// directory, keeping the recorded transcript, turns and cost, and rewrites
// summary.md / summary.json.
func rescoreRun(rootAbs, runDir string) {
	runAbs, _ := filepath.Abs(runDir)
	entries, err := os.ReadDir(runAbs)
	must(err)
	var rs []taskRun
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		ws := filepath.Join(runAbs, e.Name())
		b, err := os.ReadFile(filepath.Join(ws, "result.json"))
		if err != nil {
			continue
		}
		var tr taskRun
		must(json.Unmarshal(b, &tr))
		taskDir := filepath.Join(rootAbs, "bench", "tasks", tr.Task)
		res, err := score.Score(rootAbs, taskDir, filepath.Join(ws, "candidate"), filepath.Join(ws, "_scorer"))
		if err != nil {
			tr.Error = "rescore: " + err.Error()
		} else {
			tr.Score = res
		}
		saveJSON(filepath.Join(ws, "result.json"), tr)
		fmt.Printf("%-28s %s\n", tr.Task, oneLine(tr))
		rs = append(rs, tr)
	}
	// Recover model/tools labels from the previous summary if present.
	if b, err := os.ReadFile(filepath.Join(runAbs, "summary.md")); err == nil {
		for _, line := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(line, "model=") {
				fmt.Sscanf(line, "model=%s tools=%s", model, tools)
			}
		}
	}
	writeSummary(runAbs, rs)
}
