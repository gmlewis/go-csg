// repl implements a language repl.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/gmlewis/go-monkey/repl"
)

var (
	verbose = flag.Bool("v", false, "Verbose logging")
)

func main() {
	flag.Parse()

	user, err := user.Current()
	check("user.Current: %v", err)

	fmt.Printf("Hello %v! This is the Monkey programming language.\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
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
