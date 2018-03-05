package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// will be set by goreleaser to release tag when building
	version = "master"

	// command options
	v = flag.Bool("v", false, "print the current version")
)

func main() {
	flag.Parse()

	if *v {
		fmt.Println(version)
		os.Exit(0)
	}
}
