# go-jsontracing
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tfaller/go-jsontracing)](https://pkg.go.dev/github.com/tfaller/go-jsontracing)

jsontracing is basically a code generator to trace used JSON properties.
The generator is intended to be used with `go generate`. The generator builds special types 
to trace which properties of a JSON document got read. 

## Code generator
Let's have the following go JSON structure (file person.go).

```go
package person

type Person struct {
    Name    string
    Parents []Person `json:"parents"`
}

//go:generate go run github.com/tfaller/go-jsontracing/cmd/generator -o tracer.go -pkg person -pkgPath "my/nice/person" -t Person
```
The generator options

**-o**: The target source file, here tracer.go

**-pkg**: The package name, here person

**-pkgPath**: The full package path of this go module, here my/nice/person

**-t**: The JSON struct type that should be traced, here Person (multiple "-t" options are allowed)

Run it with `go generate person.go`.

## Tracing
Now use the generated tracer
```go
var aPersonJSON Person
json.Unmarshal([]byte(`{"Name": "Thomas", "parents": [{"Name": "foo"}, {"Name": "bar"}]}`), &aPersonJSON)

tracer := jsontracing.Tracer{}
aPerson := NewSPerson(&aPersonJSON, tracer)

// use now only aPerson ....
aPerson.Name()
parents := aPerson.Parents()
parents.Get(1).Name()

// now check used properties
fmt.Print(jsontracing.TracerToPath(tracer))
```
The console output is
```
[[Name] [parents] [parents 1] [parents 1 Name]]
```
## Generator limitations

* all JSON custom struct types have to be in the same package 
* slice and map elements can't be pointer

These limitations are currently in place. But these are no hard
limitations. It was easier for now to ignore these special cases.