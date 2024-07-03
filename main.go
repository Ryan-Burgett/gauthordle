package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/josephnaberhaus/gauthordle/internal/game"
	"github.com/josephnaberhaus/gauthordle/internal/git"
	escapes "github.com/snugfox/ansi-escapes"
)

const help = "A daily game where you try to guess the author of some Git commits.\n\nTo play, simply \"git checkout\" the main development branch of your repository\nand run this program with no arguments.\n\nNew games start at midnight Central Time."

func exit(err error) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%sERROR: %s%s", escapes.TextColorRed, err.Error(), escapes.TextColorWhite))
	os.Exit(1)
}

func main() {
	random := flag.Bool("random", false, "If true, play a random game instead of the daily game.")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Println(help)
		fmt.Println()

		flag.Usage()
		os.Exit(0)
	}

	if !git.IsGitInstalled() {
		exit(errors.New("git must be installed"))
	}

	if !git.IsInGitRepo() {
		exit(errors.New("must be in a git repository"))
	}

	fmt.Println("Building today's game...")

	var puzzle game.Puzzle
	var err error
	if *random {
		puzzle, err = game.BuildRandom()
	} else {
		puzzle, err = game.BuildToday()
	}
	if err != nil {
		exit(err)
	}

	err = puzzle.Run()
	if err != nil {
		exit(err)
	}
}
