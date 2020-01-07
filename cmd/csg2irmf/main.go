// csg2irmf reads a CSG file and writes out IRMF.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gmlewis/go-csg/evaluator"
	"github.com/gmlewis/go-csg/irmf"
	"github.com/gmlewis/go-csg/lexer"
	"github.com/gmlewis/go-csg/parser"
)

var (
	center  = flag.Bool("center", true, "Center the IRMF in world space.")
	verbose = flag.Bool("v", false, "Verbose logging")
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}

	log.Println("Done.")
}

func process(filename string) {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	le := lexer.New(string(buf))
	p := parser.New(le)
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		log.Fatalf("%v\n", strings.Join(errs, "\n"))
	}

	obj := evaluator.Eval(program, nil)

	shader := irmf.New(obj, *center)

	if shader.MBB == nil {
		log.Println("WARNING: CSG contains features that are not yet supported.")
		shader.MBB = &irmf.MBB{}
	}

	out := fmt.Sprintf(`/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [%v,%v,%v],
  min: [%v,%v,%v],
  units: "mm",
}*/

%v
`, shader.MBB.XMax, shader.MBB.YMax, shader.MBB.ZMax,
		shader.MBB.XMin, shader.MBB.YMin, shader.MBB.ZMin,
		shader.String())

	outFilename := strings.Replace(filename, ".csg", ".irmf", -1)
	log.Printf("Writing %v", outFilename)
	check("WriteFile(%q): %v", outFilename, ioutil.WriteFile(outFilename, []byte(out), 0644))
}

func check(fmtStr string, args ...interface{}) {
	if err := args[len(args)-1]; err != nil {
		log.Fatalf(fmtStr, args...)
	}
}

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}
