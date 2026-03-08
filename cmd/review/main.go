package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"review/pkg/copilot"
	"review/pkg/diff"
	"review/pkg/git"
	"review/pkg/review"
)

var (
	checkGitInstalled   = git.CheckInstalled
	latestCommitDiff    = git.LatestCommitDiff
	latestCommitMessage = git.LatestCommitMessage
	newRunner           = func() *review.Runner { return review.NewRunner(copilot.New("gpt-5-mini")) }
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("review", flag.ContinueOnError)
	fs.SetOutput(stderr)

	verbose := fs.Bool("verbose", false, "print the prompt before sending it to Copilot")
	limitBreak := fs.Bool("limit-break", false, "raise the diff line limit from 300 to 3000")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := checkGitInstalled(); err != nil {
		if errors.Is(err, git.ErrNotInstalled) {
			return git.ErrNotInstalled
		}
		return err
	}

	rawDiff, err := latestCommitDiff()
	if err != nil {
		return err
	}
	commitMessage, err := latestCommitMessage()
	if err != nil {
		return err
	}

	if diff.CountLines(rawDiff) == 0 {
		fmt.Fprintln(stdout, "No diff found")
		return nil
	}

	runner := newRunner()
	runner.Verbose = *verbose
	if *limitBreak {
		runner.MaxDiffLines = diff.LimitBreakMaxLines
	}
	if err := runner.Run(commitMessage, rawDiff, stdout, stderr); err != nil {
		if errors.Is(err, diff.ErrTooLarge) {
			return err
		}
		return err
	}
	return nil
}
