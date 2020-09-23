package jsontracing

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestTrace(t *testing.T) {
	tracer := Tracer{}
	tracer.Trace("a")
	tracer.Trace("b").Trace("1")

	if !reflect.DeepEqual(tracer, Tracer{
		"a": Tracer{},
		"b": Tracer{
			"1": Tracer{},
		},
	}) {
		t.Error("Expected other Traceing tree")
	}
}

func TestTracerToPath(t *testing.T) {
	tracer := Tracer{}

	tracer.Trace("a").Trace("1").Trace("~")
	tracer.Trace("a").Trace("2")
	tracer.Trace("b").Trace("1").Trace("#")
	tracer.Trace("c")
	tracer.Trace("d").Trace("1")

	expected := [][]string{
		{"a"},
		{"a", "1"},
		{"a", "1", "~"},
		{"a", "2"},
		{"b"},
		{"b", "1"},
		{"b", "1", "#"},
		{"c"},
		{"d"},
		{"d", "1"},
	}

	pathes := TracerToPath(tracer)
	sortPaths(pathes)

	if !reflect.DeepEqual(pathes, expected) {
		t.Errorf("Expected %v but got %v", expected, pathes)
	}
}

func sortPaths(p [][]string) {
	sort.Slice(p, func(i, j int) bool {
		a, b := p[i], p[j]
		aL, bL := len(a), len(b)

		minL := aL
		if bL < minL {
			minL = bL
		}

		for i := 0; i < minL; i++ {
			if cmp := strings.Compare(a[i], b[i]); cmp != 0 {
				return cmp < 0
			}
		}

		return aL < bL
	})
}
