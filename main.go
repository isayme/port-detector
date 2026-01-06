package main

import (
	"github.com/isayme/port-detector/src/langs"
	"github.com/isayme/port-detector/src/view"
)

func main() {
	langs.Setup()

	view.Render()
}
