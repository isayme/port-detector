package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/isayme/port-detector/src/langs"
	"github.com/isayme/port-detector/src/view"
)

var version = "dev"

func main() {
	ver := flag.Bool("version", false, "print version")
	flag.BoolVar(ver, "v", false, "print version (shorthand)")
	flag.Parse()

	if *ver {
		fmt.Println("port-detector", version)
		os.Exit(0)
	}

	langs.Setup()

	view.Render()
}
