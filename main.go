// Chrono estimates the time spent on a project by
// inspecting the git log.
//
//	% chrono
//	24 days
//	% chrono -author=some.guy@acme.com
//	12 days
//	%
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	opts := logOpts{
		author: *author,
	}

	commits, err := log(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(commits)
}

type logOpts struct {
	author string
}

type commit struct {
	date time.Time
	line string // the commit line
}

// Call git log with custom format and parse into commits
func log(opts logOpts) ([]commit, error) {
	var commits []commit

	cmd := "git"
	args := []string{"log", "--format='%H %ct'"}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return commits, err
	}

	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		commits = append(commits, commit{line: line})
	}

	return commits, nil
}
