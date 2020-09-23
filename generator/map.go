package generator

import (
	"errors"
	"reflect"
	"text/template"
)

// ErrMapIsNotAMap indicates that not a map was supplied to the function
var ErrMapIsNotAMap = errors.New("value is not a map")

// ErrMapKeyNotString indicates that the map key is not string
var ErrMapKeyNotString = errors.New("map key is not string")

// ErrMapElementIsPtr indicates that the map value is a pointer
var ErrMapElementIsPtr = errors.New("map element can't be a pointer")

type mapType struct {
	Name        string
	Type        string
	ReturnType  string
	Complex     bool
	ComplexType string
}

func writeMap(t reflect.Type, g *Generator) error {
	if t.Kind() != reflect.Map {
		return ErrMapIsNotAMap
	}

	if t.Key().Kind() != reflect.String {
		return ErrMapKeyNotString
	}

	elementType := t.Elem()
	elementKind := elementType.Kind()

	if elementKind == reflect.Ptr {
		return ErrMapElementIsPtr
	}

	g.AddType(elementType)

	mapTracerTemplate.Execute(g.out, mapType{
		Name:        typeName(t),
		Type:        typeDeclaration(elementType),
		Complex:     isComplex(elementType),
		ReturnType:  returnType(elementType),
		ComplexType: typeName(elementType),
	})

	return nil
}

var mapTracerTemplate, _ = template.New("map").Parse(`
// {{.Name}} wraps a map[string]{{.Type}}.
// This wrapper is used to trace which values of the map got used
type {{.Name}} struct {
	t jsontracing.Tracer
	m map[string]{{.Type}}
}

// Len gets the len of the underlaying map
func (m *{{.Name}}) Len() int {
	return len(m.m)
}

// Get returns a the value at the given index of the underlaying slice
func (m *{{.Name}}) Get(key string) ({{.ReturnType}}, bool) {
	v, exists := m.m[key]
	
	{{- if .Complex}}
	return New{{.ComplexType}}(&v, m.t.Trace(fmt.Sprintf("%v", key))), exists
	{{else}}
	m.t.Trace(fmt.Sprintf("%v", key))
	return v, exists
	{{end}}
}

// Range is used to iterate over all values of the underlaying slice
func (m *{{.Name}}) Range(r func (string, {{.ReturnType}}) bool) {
	for k := range m.m {
		traced, _ := m.Get(k)
		if !r(k, traced) {
			return
		}
	}
}

// New{{.Name}} creates a new tracing wrapper for map[sting]{{.Type}}
func New{{.Name}}(s *map[string]{{.Type}}, t jsontracing.Tracer) *{{.Name}} {
	return &{{.Name}}{t, *s}
}
`)
