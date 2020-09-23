package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

type typesFlag []string

func (t *typesFlag) String() string {
	return strings.Join(*t, ",")
}

func (t *typesFlag) Set(value string) error {
	*t = append(*t, value)
	return nil
}

var pkgName = flag.String("pkg", "mypackage", "package name of the auto trace types")
var pkgPath = flag.String("pkgPath", "", "import path of the package")
var outFile = flag.String("o", "", "out file")

func main() {
	var types = typesFlag{}
	flag.Var(&types, "t", "type that should be auto-traced")
	flag.Parse()

	if *outFile == "" || len(types) == 0 || *pkgPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// first ... write file that generates the dynamic generator
	generatorMain, err := template.New("").Parse(generatorGoSrc)
	if err != nil {
		log.Fatal(err)
	}

	// create tmp dir to store main file
	tmpGoDir := path.Dir(*outFile) + "/go-jsontracing-cmd"
	tmpGoFile := tmpGoDir + "/generator.go"

	err = os.Mkdir(tmpGoDir, 0744)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	f, err := os.Create(tmpGoFile)
	if err != nil {
		log.Fatal(err)
	}

	err = generatorMain.Execute(f, struct {
		Type    []string
		PkgPath string
		Pkg     string
		OutFile string
	}{
		Type:    types,
		Pkg:     *pkgName,
		PkgPath: *pkgPath,
		OutFile: *outFile,
	})

	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// run that file ... which actually generates the types
	cmd := exec.Command("go", "run", tmpGoFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		// exit with fatal.
		// this will not cleanup the files ... this is actually nice
		// because this gives the dev a chance to trace a problem in the
		// auto generated source files
		log.Fatal(err)
	}

	// cleanup ...
	err = os.Remove(tmpGoFile)
	if err != nil {
		log.Print(err)
	}
	err = os.Remove(tmpGoDir)
	if err != nil {
		log.Print(err)
	}

	// last step ... format the code
	cmd = exec.Command("go", "fmt", *outFile)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// generatorGoSrc is a little dynamic generated go
// program to build the auto-tracing types
const generatorGoSrc = `
package main

import (
	"log"
	"os"
	"reflect"

	"github.com/tfaller/go-jsontracing/generator"
	"{{.PkgPath}}"
)

func main() {
	f, err := os.Create("{{.OutFile}}")
	if err != nil {
		log.Fatal(err)
	}

	g := generator.NewGenerator(f, "{{.Pkg}}")

	{{range .Type}}
	g.AddType(reflect.TypeOf((*{{$.Pkg}}.{{.}})(nil)).Elem())
	{{end}}

	err = g.Generate()
	if err != nil {
		log.Fatal(err)
	}
}`
