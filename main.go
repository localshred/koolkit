package main

import (
	"os"

	"github.com/localshred/koolkit/cmd"
)

// version is overwritten at release time via -ldflags.
var version = "dev"

func main() { os.Exit(cmd.Execute(version)) }
