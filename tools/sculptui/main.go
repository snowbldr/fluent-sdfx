// Command sculptui is the first test loop for sculpting fluent-sdfx solids
// from a browser: a 3D view of a form, click and drag on it to pull the
// surface, every pull a step in a history file on disk. Edit the file by
// hand and the view reloads; turn a step off or delete it and the form
// reflows as if it never happened.
//
//	go run ./tools/sculptui              # history in ./sculpt.yaml
//	go run ./tools/sculptui -file nose.yaml -addr 127.0.0.1:7175
//
// The file is a fluent-sdfx history (package history): any op the library
// has can be written into it by hand or by an agent, and the view shows
// it. The UI itself makes pulls, and edits the numbers of a pull.
package main

import (
	"bytes"
	"embed"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf16"

	"github.com/snowbldr/fluent-sdfx/history"
	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
)

//go:embed web
var webFS embed.FS

type server struct {
	path         string
	cells        int // final mesh resolution, cells on the longest side
	previewCells int // while dragging or editing

	// rt caches every step's result by what it is and what it depends on,
	// so a preview of an edit recomputes only the edited step and those
	// after it, and undo finds the old results still there.
	rt *history.Runtime

	mu      sync.Mutex
	hist    *history.History
	version int
	modTime time.Time // of the file as we last wrote or read it
	loadErr string    // the file on disk does not parse; the last good history is shown

	meshMu      sync.Mutex
	meshVersion int
	meshBytes   []byte
	meshInfo    meshInfo
}

type meshInfo struct {
	Triangles int          `json:"triangles"`
	Millis    int64        `json:"millis"`
	Steps     []stepStatus `json:"steps,omitempty"` // the steps that are not ok
	Error     string       `json:"error,omitempty"`
}

type stepStatus struct {
	Index  int            `json:"index"`
	Status history.Status `json:"status"`
	Error  string         `json:"error,omitempty"`
}

func main() {
	path := flag.String("file", "sculpt.yaml", "history file; created if missing")
	addr := flag.String("addr", "127.0.0.1:7175", "listen address")
	cells := flag.Int("cells", 140, "mesh resolution after a change (cells on the longest side)")
	preview := flag.Int("preview", 90, "mesh resolution while dragging")
	bakeEvery := flag.Int("bake-every", 8, "sample the form onto a grid after this many pulls in a row (0 never)")
	bakeCells := flag.Int("bake", 128, "grid resolution of a bake")
	flag.Parse()

	s := &server{path: *path, cells: *cells, previewCells: *preview, meshVersion: -1,
		rt: history.NewRuntime(history.Options{BakeEvery: *bakeEvery, BakeCells: *bakeCells})}
	if h, err := history.Load(*path); err == nil {
		s.hist = h
	} else if old, lerr := legacy(*path); lerr == nil {
		// A file from the first version of this tool: keep a copy, rewrite it.
		if err := os.Rename(*path, *path+".v1"); err != nil {
			log.Fatal(err)
		}
		if err := old.Save(*path); err != nil {
			log.Fatal(err)
		}
		log.Printf("converted %s to a history; the old file is %s.v1", *path, *path)
		s.hist = old
	} else if errors.Is(err, os.ErrNotExist) {
		s.hist, _ = starter("figure")
		if err := s.hist.Save(*path); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	if fi, err := os.Stat(*path); err == nil {
		s.modTime = fi.ModTime()
	}
	go s.watch()

	web, _ := fs.Sub(webFS, "web")
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.FS(web)))
	mux.HandleFunc("GET /api/state", s.handleState)
	mux.HandleFunc("GET /api/mesh", s.handleMesh)
	mux.HandleFunc("POST /api/preview", s.handlePreview)
	mux.HandleFunc("POST /api/steps", s.handleAdd)
	mux.HandleFunc("PUT /api/steps/{i}", s.handleEdit)
	mux.HandleFunc("POST /api/steps/{i}/toggle", s.handleToggle)
	mux.HandleFunc("DELETE /api/steps/{i}", s.handleDelete)
	mux.HandleFunc("POST /api/undo", s.handleUndo)
	mux.HandleFunc("POST /api/start", s.handleStart)

	log.Printf("sculptui: http://%s  (history: %s)", *addr, *path)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

// watch reloads the history when the file changes under us: a hand edit,
// or an agent editing it. A file that does not parse leaves the last good
// history on screen and reports the error.
func (s *server) watch() {
	for range time.Tick(400 * time.Millisecond) {
		fi, err := os.Stat(s.path)
		if err != nil {
			continue
		}
		s.mu.Lock()
		changed := !fi.ModTime().Equal(s.modTime)
		s.mu.Unlock()
		if !changed {
			continue
		}
		h, err := history.Load(s.path)
		s.mu.Lock()
		s.modTime = fi.ModTime()
		if err != nil {
			s.loadErr = err.Error()
		} else {
			s.loadErr = ""
			s.hist = h
		}
		s.version++
		s.mu.Unlock()
		log.Printf("reloaded %s (version %d)", s.path, s.version)
	}
}

// commit applies a change to a copy of the history, writes the file and
// bumps the version. The change runs under the lock; a change that fails,
// or leaves the history invalid, leaves everything as it was.
func (s *server) commit(w http.ResponseWriter, change func(h *history.History) error) {
	s.mu.Lock()
	h := s.hist.Clone()
	if err := change(h); err != nil {
		s.mu.Unlock()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Save(s.path); err != nil {
		s.mu.Unlock()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if fi, err := os.Stat(s.path); err == nil {
		s.modTime = fi.ModTime()
	}
	s.hist = h
	s.loadErr = ""
	s.version++
	s.mu.Unlock()
	s.handleState(w, nil)
}

func (s *server) handleState(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, map[string]any{
		"version": s.version,
		"history": s.hist,
		"file":    s.path,
		"loadErr": s.loadErr,
	})
}

// handleMesh returns the current history's mesh: little-endian float32
// triangle corners, nine per triangle. Build details ride in a header.
func (s *server) handleMesh(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	h, v := s.hist, s.version
	s.mu.Unlock()

	s.meshMu.Lock()
	defer s.meshMu.Unlock()
	if s.meshVersion != v {
		s.meshBytes, s.meshInfo = s.meshOf(h, s.cells)
		s.meshVersion = v
	}
	writeMesh(w, s.meshBytes, s.meshInfo, v)
}

// handlePreview meshes an unsaved variant of the history at preview
// resolution: the live view while a drag is under way or a step's form is
// being edited. The body is the full step list to show; the runtime's
// cache means only the steps that differ, and those after them, are
// evaluated.
func (s *server) handlePreview(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Steps []history.Step `json:"steps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	h := s.hist.Clone()
	v := s.version
	s.mu.Unlock()
	h.Steps = req.Steps
	if err := h.Validate(); err != nil {
		writeMesh(w, nil, meshInfo{Error: err.Error()}, v)
		return
	}
	b, info := s.meshOf(h, s.previewCells)
	writeMesh(w, b, info, v)
}

func stepIndex(r *http.Request, h *history.History) (int, error) {
	i, err := strconv.Atoi(r.PathValue("i"))
	if err != nil || i < 0 || i >= len(h.Steps) {
		return 0, fmt.Errorf("no step %s", r.PathValue("i"))
	}
	return i, nil
}

// handleEdit replaces one step with the one in the body.
func (s *server) handleEdit(w http.ResponseWriter, r *http.Request) {
	var st history.Step
	if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.commit(w, func(h *history.History) error {
		i, err := stepIndex(r, h)
		if err != nil {
			return err
		}
		return h.Set(i, st)
	})
}

func (s *server) handleAdd(w http.ResponseWriter, r *http.Request) {
	var add []history.Step
	if err := json.NewDecoder(r.Body).Decode(&add); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.commit(w, func(h *history.History) error {
		for _, st := range add {
			if err := h.Insert(len(h.Steps), st); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *server) handleToggle(w http.ResponseWriter, r *http.Request) {
	s.commit(w, func(h *history.History) error {
		i, err := stepIndex(r, h)
		if err != nil {
			return err
		}
		return h.SetOff(i, !h.Steps[i].Off)
	})
}

func (s *server) handleDelete(w http.ResponseWriter, r *http.Request) {
	s.commit(w, func(h *history.History) error {
		i, err := stepIndex(r, h)
		if err != nil {
			return err
		}
		return h.Remove(i)
	})
}

func (s *server) handleUndo(w http.ResponseWriter, _ *http.Request) {
	s.commit(w, func(h *history.History) error {
		if len(h.Steps) == 0 {
			return errors.New("nothing to undo")
		}
		return h.Remove(len(h.Steps) - 1)
	})
}

// handleStart replaces the whole history with a starter form.
func (s *server) handleStart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Form string `json:"form"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.commit(w, func(h *history.History) error {
		fresh, err := starter(req.Form)
		if err != nil {
			return err
		}
		*h = *fresh
		return nil
	})
}

// meshOf evaluates the history and meshes its bodies with cells on the
// longest side of their bounds.
func (s *server) meshOf(h *history.History, cells int) ([]byte, meshInfo) {
	t0 := time.Now()
	res := s.rt.Eval(h)
	var info meshInfo
	for _, st := range res.Steps {
		if st.Status == history.Failed || st.Status == history.Blocked {
			info.Steps = append(info.Steps, stepStatus{Index: st.Index, Status: st.Status, Error: st.Error})
		}
	}
	sol := res.Solid()
	if sol == nil {
		info.Error = "the history makes no solid"
		return nil, info
	}
	tris := mesh.CollectTriangles(sol, render.NewMarchingCubesOctreeParallel(cells))
	b := make([]byte, 0, len(tris)*36)
	for _, t := range tris {
		for _, p := range t {
			b = binary.LittleEndian.AppendUint32(b, math.Float32bits(float32(p.X)))
			b = binary.LittleEndian.AppendUint32(b, math.Float32bits(float32(p.Y)))
			b = binary.LittleEndian.AppendUint32(b, math.Float32bits(float32(p.Z)))
		}
	}
	info.Triangles, info.Millis = len(tris), time.Since(t0).Milliseconds()
	return b, info
}

func writeMesh(w http.ResponseWriter, b []byte, info meshInfo, version int) {
	j, _ := json.Marshal(info)
	w.Header().Set("X-Mesh-Info", asciiJSON(j))
	w.Header().Set("X-Version", strconv.Itoa(version))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(b)
}

// asciiJSON escapes every non-ASCII rune as \uXXXX: still the same JSON,
// but safe in an HTTP header, which browsers read as Latin-1.
func asciiJSON(j []byte) string {
	var b strings.Builder
	for _, r := range string(j) {
		if r < 128 {
			b.WriteRune(r)
			continue
		}
		for _, u := range utf16.Encode([]rune{r}) {
			fmt.Fprintf(&b, "\\u%04x", u)
		}
	}
	return b.String()
}

func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
