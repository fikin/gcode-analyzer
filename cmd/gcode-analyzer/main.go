package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fikin/gcode-analyzer/gcode"
)

var doc = struct {
	Name string
	Doc  string
}{
	Name: "gcode-analyzer",
	Doc: `scans gcode files and analyzes its content.
`,
}

var gCodeFile string

func main() {
	addCmdlineFlags()

	flag.Parse() // (ExitOnError)

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	err := run(args[0])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func addCmdlineFlags() {
	flag.Usage = func() {
		paras := strings.Split(doc.Doc, "\n\n")
		fmt.Fprintf(os.Stderr, "%s: %s\n\n", doc.Name, paras[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [-flag] [gcode file]\n\n", doc.Name)
		if len(paras) > 1 {
			fmt.Fprintln(os.Stderr, strings.Join(paras[1:], "\n\n"))
		}
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}
}

func run(filename string) error {
	content, err := gcode.ParseGCodeFile(filename)
	if err != nil {
		return err
	}
	gcode.MarshallAsText(content, os.Stdout)
	return nil
}
