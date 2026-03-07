package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"review/pkg/copilot"
	"review/pkg/diff"
	"review/pkg/git"
	"review/pkg/review"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("review", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	verbose := fs.Bool("verbose", false, "print the prompt before sending it to Copilot")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := git.CheckInstalled(); err != nil {
		if errors.Is(err, git.ErrNotInstalled) {
			return git.ErrNotInstalled
		}
		return err
	}

	rawDiff, err := git.LatestCommitDiff()
	if err != nil {
		return err
	}

	if diff.CountLines(rawDiff) == 0 {
		fmt.Fprintln(os.Stdout, "No diff found")
		return nil
	}

	runner := review.NewRunner(copilot.New("gpt-5-mini"))
	runner.Verbose = *verbose
	if err := runner.Run(rawDiff, os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, diff.ErrTooLarge) {
			return err
		}
		return err
	}
	return nil
}
