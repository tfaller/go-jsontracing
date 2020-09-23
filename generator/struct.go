package generator

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

type structType struct {
	Name   string
	Type   string
	Fields []structJSONField
}

type structJSONField struct {
	Name       string
	Type       string
	ReturnType string
	JSON       string
	Complex    bool
	Ptr        bool
}

func writeStruct(s reflect.Type, g *Generator) error {
	if structErr != nil {
		panic(structErr)
	}

	if s.Kind() != reflect.Struct {
		return fmt.Errorf("type %q is not a struct", s)
	}

	fields := []structJSONField{}

	for i := s.NumField() - 1; i >= 0; i-- {
		field := s.Field(i)
		fieldType := field.Type

		isPtr := fieldType.Kind() == reflect.Ptr
		if isPtr {
			fieldType = fieldType.Elem()
		}

		g.AddType(fieldType)

		complex := isComplex(fieldType)

		fields = append(fields, structJSONField{
			JSON:       jsonFieldName(field),
			Name:       field.Name,
			Type:       typeName(fieldType),
			ReturnType: returnType(field.Type),
			Complex:    complex,
			Ptr:        isPtr,
		})
	}

	return structTemplate.Execute(g.out, structType{Name: typeName(s), Fields: fields, Type: s.Name()})
}

func jsonFieldName(field reflect.StructField) string {
	name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	if name == "" {
		name = field.Name
	}
	return name
}

var structTemplate, structErr = template.New("struct").Parse(`
// {{.Name}} wraps a {{.Type}}.
// This wrapper is used to trace which fields of the struct got used
type {{.Name}} struct {
	t jsontracing.Tracer
	v *{{.Type}}
}

{{range .Fields -}}

// {{.Name}} gets the value of the {{.Name}} field.
func (s *{{$.Name}}) {{.Name}}() {{.ReturnType}} {
	v := s.v.{{.Name}}

	{{- if .Complex}}
	return New{{.Type}}({{if not .Ptr}}&{{end}}v, s.t.Trace("{{.JSON}}")) 	

	{{- else}}
	s.t.Trace("{{.JSON}}")
	return v

	{{- end}} 
}
{{end}}

// New{{.Name}} creates a new tracer for the given value
func New{{.Name}}(v *{{.Type}}, t jsontracing.Tracer) *{{.Name}} {
	return &{{.Name}}{t, v}
}
`)
