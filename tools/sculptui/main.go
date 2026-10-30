// Command sculptui is the first test loop for sculpting fluent-sdfx solids
// from a browser: a 3D view of a form, click and drag on it to pull the
// surface, every pull a step in a YAML history on disk. Edit the file by
// hand and the view reloads; turn a step off or delete it and the form
// reflows as if it never happened.
//
//	go run ./tools/sculptui              # history in ./sculpt.yaml
//	go run ./tools/sculptui -file nose.yaml -addr 127.0.0.1:7175
//
// It exists to find out whether this way of working feels right before the
// full history runtime in meta/sculpt-cad-build-plan.md is built. It knows
// one operation (pull) and three base forms.
package main

import (
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

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
)

//go:embed web
var webFS embed.FS

type server struct {
	path         string
	cells        int // final mesh resolution, cells on the longest side
	previewCells int // while dragging

	mu      sync.Mutex
	hist    History
	version int
	modTime time.Time // of the file as we last wrote or read it
	loadErr string    // the file on disk does not parse; the last good history is shown

	meshMu      sync.Mutex
	meshVersion int
	meshBytes   []byte
	meshInfo    meshInfo
}

type meshInfo struct {
	Triangles int        `json:"triangles"`
	Millis    int64      `json:"millis"`
	StepError *StepError `json:"stepError,omitempty"`
	Error     string     `json:"error,omitempty"`
}

func main() {
	path := flag.String("file", "sculpt.yaml", "history file; created if missing")
	addr := flag.String("addr", "127.0.0.1:7175", "listen address")
	cells := flag.Int("cells", 140, "mesh resolution after a change (cells on the longest side)")
	preview := flag.Int("preview", 90, "mesh resolution while dragging")
	flag.Parse()

	s := &server{path: *path, cells: *cells, previewCells: *preview, meshVersion: -1}
	if h, err := loadHistory(*path); err == nil {
		s.hist = h
	} else if errors.Is(err, os.ErrNotExist) {
		s.hist = newHistory()
		if err := s.hist.save(*path); err != nil {
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
	mux.HandleFunc("POST /api/steps/{i}/toggle", s.handleToggle)
	mux.HandleFunc("DELETE /api/steps/{i}", s.handleDelete)
	mux.HandleFunc("POST /api/undo", s.handleUndo)
	mux.HandleFunc("POST /api/base", s.handleBase)

	log.Printf("sculptui: http://%s  (history: %s)", *addr, *path)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

// watch reloads the history when the file changes under us: a hand edit,
// or an AI editing it. A file that does not parse leaves the last good
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
		h, err := loadHistory(s.path)
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

// commit applies a change to the history, writes the file and bumps the
// version. The change runs under the lock.
func (s *server) commit(w http.ResponseWriter, change func(h *History) error) {
	s.mu.Lock()
	h := s.hist
	h.Steps = append([]Step(nil), s.hist.Steps...)
	if err := change(&h); err != nil {
		s.mu.Unlock()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.save(s.path); err != nil {
		s.mu.Unlock()
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		s.meshBytes, s.meshInfo = meshOf(h, s.cells)
		s.meshVersion = v
	}
	writeMesh(w, s.meshBytes, s.meshInfo, v)
}

// handlePreview builds the history plus the steps in flight, unsaved, at
// preview resolution: the live view while a drag is under way.
func (s *server) handlePreview(w http.ResponseWriter, r *http.Request) {
	var extra []Step
	if err := json.NewDecoder(r.Body).Decode(&extra); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	h := s.hist
	h.Steps = append(append([]Step(nil), s.hist.Steps...), extra...)
	v := s.version
	s.mu.Unlock()
	b, info := meshOf(h, s.previewCells)
	writeMesh(w, b, info, v)
}

func (s *server) handleAdd(w http.ResponseWriter, r *http.Request) {
	var add []Step
	if err := json.NewDecoder(r.Body).Decode(&add); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.commit(w, func(h *History) error {
		h.Steps = append(h.Steps, add...)
		return nil
	})
}

func (s *server) handleToggle(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("i"))
	s.commit(w, func(h *History) error {
		if err != nil || i < 0 || i >= len(h.Steps) {
			return fmt.Errorf("no step %s", r.PathValue("i"))
		}
		h.Steps[i].Off = !h.Steps[i].Off
		return nil
	})
}

func (s *server) handleDelete(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("i"))
	s.commit(w, func(h *History) error {
		if err != nil || i < 0 || i >= len(h.Steps) {
			return fmt.Errorf("no step %s", r.PathValue("i"))
		}
		h.Steps = append(h.Steps[:i], h.Steps[i+1:]...)
		return nil
	})
}

func (s *server) handleUndo(w http.ResponseWriter, _ *http.Request) {
	s.commit(w, func(h *History) error {
		if len(h.Steps) == 0 {
			return errors.New("nothing to undo")
		}
		h.Steps = h.Steps[:len(h.Steps)-1]
		return nil
	})
}

func (s *server) handleBase(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Base  string `json:"base"`
		Clear bool   `json:"clear"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.commit(w, func(h *History) error {
		if _, err := base(req.Base); err != nil {
			return err
		}
		h.Base = req.Base
		if req.Clear {
			h.Steps = nil
		}
		return nil
	})
}

// meshOf builds the history and meshes the result with cells on the
// longest side of its bounds.
func meshOf(h History, cells int) ([]byte, meshInfo) {
	t0 := time.Now()
	sol, stepErr, err := h.build()
	if err != nil {
		return nil, meshInfo{Error: err.Error()}
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
	return b, meshInfo{Triangles: len(tris), Millis: time.Since(t0).Milliseconds(), StepError: stepErr}
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

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
