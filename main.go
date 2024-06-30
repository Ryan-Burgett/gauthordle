package main

import (
	"errors"
	"fmt"
	"github.com/JosephNaberhaus/gauthordle/internal/game"
	"github.com/JosephNaberhaus/gauthordle/internal/git"
	escapes "github.com/snugfox/ansi-escapes"
	"os"
)

func exit(err error) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%sERROR: %s%s", escapes.TextColorRed, err.Error(), escapes.TextColorWhite))
	os.Exit(1)
}

func main() {
	if !git.IsGitInstalled() {
		exit(errors.New("git must be installed"))
	}

	if !git.IsInGitRepo() {
		exit(errors.New("must be in a git repository"))
	}

	println("Building today's game...")
	puzzle, err := game.BuildToday()
	if err != nil {
		exit(err)
	}

	err = puzzle.Run()
	if err != nil {
		exit(err)
	}
}
