package jsontracing

// Tracer traces what properties of a JSON document got used.
type Tracer map[string]Tracer

// Trace traces what given JSON object fields or
// array indices was read. It also returns another Tracer
// so that sub properties of an filed or index can be traced,
// if the value itself was an object or array.
func (t Tracer) Trace(name string) Tracer {
	entry := t[name]
	if entry == nil {
		// entry was not yet read
		// create a sub Tracer if the property
		// value is a object or array
		entry = Tracer{}
		t[name] = entry
	}
	return entry
}

// TracerToPath flattens all read path properties into a slice of
// of all met property paths.
func TracerToPath(t Tracer) [][]string {
	return TracerToPathWithPath(nil, t)
}

// TracerToPathWithPath is basically the same as TracerToPath but
// allows to set a path prefix which all paths of the list will have.
func TracerToPathWithPath(path []string, tracer Tracer) [][]string {
	paths := [][]string{}

	for k, v := range tracer {
		p := append(path, k)
		// add own path
		paths = append(paths, p)
		// add sub paths
		paths = append(paths, TracerToPathWithPath(p, v)...)
	}

	return paths
}
