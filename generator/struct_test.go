package generator

import (
	"reflect"
	"testing"
)

func TestJsonFieldName(t *testing.T) {
	structType := reflect.TypeOf(struct {
		Name    string
		Parents []string `json:"parents"`
		Age     int      `json:"age,omitempty"`
	}{})

	expected := []string{"Name", "parents", "age"}

	for i, name := range expected {
		if jsonFieldName(structType.Field(i)) != name {
			t.Errorf("%v: Wrong name", i)
		}
	}
}
