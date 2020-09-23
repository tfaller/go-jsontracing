package generator

import (
	"reflect"
	"strings"
)

// isComplex tells whether the type needs a custom tracer
func isComplex(t reflect.Type) bool {
	kind := t.Kind()
	return kind == reflect.Struct ||
		kind == reflect.Map ||
		kind == reflect.Slice ||
		kind == reflect.Interface
}

// typeName returns how the wrapper type is called
func typeName(t reflect.Type) string {
	switch t.Kind() {

	case reflect.Struct:
		return "S" + t.Name()

	case reflect.Map:
		return "Map" + strings.Title(typeName(t.Elem()))

	case reflect.Slice:
		return "Slice" + strings.Title(typeName(t.Elem()))

	case reflect.Interface:
		if t.NumMethod() != 0 {
			panic("only empty interfaces are allowed")
		}
		return "Interface"
	}
	return t.Name()
}

// typeDeclaration returns how a type would
// be declared as variable in a go source file
func typeDeclaration(t reflect.Type) string {
	switch t.Kind() {

	case reflect.Map:
		return "map[" + typeDeclaration(t.Key()) + "]" + typeDeclaration(t.Elem())

	case reflect.Slice:
		return "[]" + typeDeclaration(t.Elem())

	case reflect.Interface:
		if t.NumMethod() != 0 {
			panic("only empty interfaces are allowed")
		}
		return "interface{}"
	}

	// struct, int, float, string, boolean ...
	return t.Name()
}

// returnType converts a reflect type to its
// correct tracer function return type.
func returnType(t reflect.Type) string {
	kind := t.Kind()
	if kind == reflect.Interface {
		// interfaces stays a interface
		return "interface{}"
	}

	isPtr := kind == reflect.Ptr
	if isPtr {
		t = t.Elem()
	}

	if isPtr || isComplex(t) {
		return "*" + typeName(t)
	}
	return typeName(t)
}
