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
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	// version will be set by goreleaser to release tag when building
	v = "master"

	// command options
	author  = flag.String("author", "", "the commit author")
	skip    = flag.String("skip", "60m", "skip commits with larger delta than this")
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

	skip, err := time.ParseDuration(*skip)
	if err != nil {
		fmt.Println("invalid value for skip")
		os.Exit(1)
	}

	commits, err := log(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var dur time.Duration

	for i, _ := range commits {
		var curr, next commit
		curr = commits[i]

		if i < len(commits)-1 {
			next = commits[i+1]
		} else {
			break
		}

		delta := curr.time.Sub(next.time)

		if delta < skip {
			dur = dur + delta
		}
	}

	fmt.Println(dur)
}

type logOpts struct {
	author string
}

type commit struct {
	time time.Time
	hash string
	raw  string // raw commit line from git log --format={...}
}

// Call git log with custom format and parse into commits
func log(opts logOpts) ([]commit, error) {
	var commits []commit

	cmd := "git"
	args := []string{"log", "--format=%H %ct"}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return commits, err
	}

	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		c, err := parseCommitLine(line)
		if err != nil {
			continue
		}
		commits = append(commits, c)
	}

	return commits, nil
}

// Parses a line from custom format into a commit struct
func parseCommitLine(line string) (commit, error) {
	c := commit{
		raw: line,
	}
	parts := strings.Split(line, " ")

	if len(parts) != 2 {
		return c, errors.New("commit line could not be parsed: unexpected format")
	}

	c.hash = parts[0]

	i, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return c, err
	}

	c.time = time.Unix(i, 0)

	return c, nil
}
