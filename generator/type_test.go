package generator

import (
	"reflect"
	"testing"
)

func TestTypeName(t *testing.T) {
	testCases := []struct {
		t    reflect.Type
		name string
	}{
		{
			reflect.TypeOf((*interface{})(nil)).Elem(),
			"Interface",
		},
		{
			reflect.TypeOf((*int)(nil)).Elem(),
			"int",
		},
	}

	for i, test := range testCases {
		if n := typeName(test.t); n != test.name {
			t.Errorf("%v: Expected %q but got %q", i, test.name, n)
		}
	}
}

func TestTypeDeclaration(t *testing.T) {
	testCases := []struct {
		t    reflect.Type
		name string
	}{
		{
			reflect.TypeOf((*interface{})(nil)).Elem(),
			"interface{}",
		},
		{
			reflect.TypeOf((*int)(nil)).Elem(),
			"int",
		},
		{
			reflect.TypeOf((*map[string]interface{})(nil)).Elem(),
			"map[string]interface{}",
		},
		{
			reflect.TypeOf((*[]string)(nil)).Elem(),
			"[]string",
		},
	}

	for i, test := range testCases {
		if n := typeDeclaration(test.t); n != test.name {
			t.Errorf("%v: Expected %q but got %q", i, test.name, n)
		}
	}
}

func TestTypeDeclarationInterface(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should have paniced")
		}
	}()

	typeDeclaration(reflect.TypeOf((*interface{ Method() })(nil)).Elem())
}
