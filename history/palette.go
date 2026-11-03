package history

import (
	"reflect"
	"sort"
	"strings"
)

// OpInfo describes one op for a UI palette or an agent's tool list: what
// to write in `op`, what it does, and what its args are. It comes from the
// same source as the evaluator, so it cannot describe an op that does not
// run or miss one that does.
type OpInfo struct {
	Op string `json:"op"` // as written in a step: solid.Sphere, or Translate for a method
	// Kind is function or method. A method's input is the step named in
	// `on`, or the previous step.
	Kind       string      `json:"kind"`
	Receiver   string      `json:"receiver,omitempty"` // a method's input type, solid.Solid
	Doc        string      `json:"doc"`
	Params     []ParamInfo `json:"params"`
	Returns    string      `json:"returns"`
	Registered bool        `json:"registered,omitempty"`
	// Input is set when a function's first parameter is a form that an
	// omitted `on` fills with the previous step.
	Input string `json:"input,omitempty"`
}

// ParamInfo is one parameter.
type ParamInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Widget   string `json:"widget"` // number, length, angle, count, toggle, text, point, offset, direction, size, step, blend, points, struct, value
	Variadic bool   `json:"variadic,omitempty"`
	// Fields are a struct parameter's fields, which args may give directly.
	Fields []ParamInfo `json:"fields,omitempty"`
}

// Palette describes every op: the library functions, the methods on every
// type a step can produce, and the registered functions, sorted by op.
func Palette() []OpInfo {
	var out []OpInfo
	for _, name := range funcNames() {
		c, _ := lookupFunc(name)
		info := describeFunc(name, c.fn.Type(), c.params, c.doc, 0)
		info.Kind, info.Registered = "function", c.reg != nil
		if t := c.fn.Type(); t.NumIn() > 0 && formType(t.In(0)) && !(t.IsVariadic() && t.NumIn() == 1) {
			info.Input = c.params[0]
		}
		out = append(out, info)
	}
	var recvs []string
	for k := range genTypes {
		recvs = append(recvs, k)
	}
	sort.Strings(recvs)
	for _, rk := range recvs {
		pt := genTypes[rk]
		for i := 0; i < pt.NumMethod(); i++ {
			m := pt.Method(i)
			s, known := genSigs[rk+"."+m.Name]
			if !known {
				continue // filtered out by opsgen (writes a file, no result, ...)
			}
			// The method's type includes the receiver as its first parameter.
			info := describeFunc(m.Name, m.Type, append([]string{"recv"}, s.params...), s.doc, 1)
			info.Kind, info.Receiver = "method", rk
			out = append(out, info)
		}
	}
	return out
}

func describeFunc(op string, t reflect.Type, params []string, doc string, skip int) OpInfo {
	info := OpInfo{Op: op, Doc: doc, Params: []ParamInfo{}}
	if t.NumOut() > 0 {
		info.Returns = typeName(t.Out(0))
	}
	for k := skip; k < t.NumIn(); k++ {
		pt := t.In(k)
		p := ParamInfo{Name: params[k]}
		if t.IsVariadic() && k == t.NumIn()-1 {
			p.Variadic = true
			pt = pt.Elem()
		}
		p.Type, p.Widget = typeName(pt), widget(p.Name, pt)
		if p.Widget == "struct" {
			p.Fields = structFields(pt)
		}
		info.Params = append(info.Params, p)
	}
	return info
}

func structFields(t reflect.Type) []ParamInfo {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	var out []ParamInfo
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		out = append(out, ParamInfo{Name: f.Name, Type: typeName(f.Type), Widget: widget(f.Name, f.Type)})
	}
	return out
}

// widget is the control a UI would show for a parameter, from its type
// and, for vectors and numbers, its name.
func widget(name string, t reflect.Type) string {
	n := strings.ToLower(name)
	has := func(words ...string) bool {
		for _, w := range words {
			if n == w || strings.Contains(n, w) {
				return true
			}
		}
		return false
	}
	switch {
	case formType(t) || t == anchoredT:
		return "step"
	case t == vec3Type || t == vec2Type:
		switch {
		case has("size", "dims"):
			return "size"
		case has("axis", "dir", "normal", "up"):
			return "direction"
		case n == "v" || n == "by" || has("delta", "offset", "step", "shift"):
			return "offset"
		}
		return "point"
	case t == vec3sType:
		return "points"
	case t.Kind() == reflect.Func:
		return "blend"
	case t.Kind() == reflect.Float64 || t.Kind() == reflect.Float32:
		if has("deg", "angle", "theta") {
			return "angle"
		}
		if has("radius", "height", "width", "depth", "thickness", "length", "round", "dist", "amount", "gap", "size") || n == "r" || n == "h" || n == "k" || n == "d" {
			return "length"
		}
		return "number"
	case t.Kind() >= reflect.Int && t.Kind() <= reflect.Uint64:
		return "count"
	case t.Kind() == reflect.Bool:
		return "toggle"
	case t.Kind() == reflect.String:
		return "text"
	case isParamStruct(t):
		return "struct"
	}
	return "value"
}
