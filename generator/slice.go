package generator

import (
	"errors"
	"reflect"
	"text/template"
)

// ErrSliceElementIsPtr indicates that the element is a pointer type
var ErrSliceElementIsPtr = errors.New("slice element can't be a pointer")

type sliceType struct {
	Name        string
	Type        string
	ReturnType  string
	Complex     bool
	ComplexType string
}

func writeSlice(t reflect.Type, g *Generator) error {
	elementType := t.Elem()

	if elementType.Kind() == reflect.Ptr {
		return ErrSliceElementIsPtr
	}

	g.AddType(elementType)

	return sliceTracerTemplate.Execute(g.out, sliceType{
		Name:        typeName(t),
		Type:        typeDeclaration(elementType),
		ReturnType:  returnType(elementType),
		Complex:     isComplex(elementType),
		ComplexType: typeName(elementType)},
	)
}

var sliceTracerTemplate, _ = template.New("slice").Parse(`
// {{.Name}} wraps a []{{.Type}}.
// This wrapper is used to trace which values of the slice got used
type {{.Name}} struct {
	t jsontracing.Tracer
	slice []{{.Type}}
}

// Len gets the len of the underlaying slice
func (s *{{.Name}}) Len() int {
	return len(s.slice)
}

// Get returns a the value at the given index of the underlaying slice
func (s *{{.Name}}) Get(idx int) {{.ReturnType}} {
	v := s.slice[idx]
	
	{{- if .Complex}}
	return New{{.ComplexType}}(&v, s.t.Trace(fmt.Sprintf("%v", idx)))
	{{else}}
	s.t.Trace(fmt.Sprintf("%v", idx))
	return v
	{{end}}
}

// Range is used to iterate over all values of the underlaying slice
func (s *{{.Name}}) Range(r func (int, {{.ReturnType}}) bool) {
	len := len(s.slice)
	for i := 0; i < len; i ++ {
		if !r(i, s.Get(i)) {
			return
		}
	}
}

// New{{.Name}} creates a new tracing wrapper for []{{.Type}}
func New{{.Name}}(s *[]{{.Type}}, t jsontracing.Tracer) *{{.Name}} {
	return &{{.Name}}{t, *s}
}
`)
