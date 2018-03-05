package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	// version will be set by goreleaser to release tag when building
	v = "master"

	// command options
	author  = flag.String("author", "", "the commit author")
	version = flag.Bool("v", false, "print the current version")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Println(v)
		os.Exit(0)
	}

	if *author == "" {
		fmt.Println("author must be provided")
		os.Exit(1)
	}

	opts := logOpts{
		author: *author,
	}

	_, err := log(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type logOpts struct {
	author string
}

type commit struct {
	date time.Time
}

// Call git log with format pretty and parse into commits
func log(opts logOpts) ([]commit, error) {
	var commits []commit
	return commits, nil
}
