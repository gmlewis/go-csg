// csg2irmf reads a CSG file and writes out IRMF.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gmlewis/go-csg/evaluator"
	"github.com/gmlewis/go-csg/lexer"
	"github.com/gmlewis/go-csg/object"
	"github.com/gmlewis/go-csg/parser"
)

var (
	verbose = flag.Bool("v", false, "Verbose logging")
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}
}

func process(filename string) {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	env := object.NewEnvironment()
	le := lexer.New(string(buf))
	p := parser.New(le)
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		fmt.Fprintf(os.Stderr, "%v\n", strings.Join(errs, "\n"))
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Printf("%v\n", evaluated.Inspect())
	}
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
