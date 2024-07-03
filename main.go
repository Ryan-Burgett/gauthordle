package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/josephnaberhaus/gauthordle/internal/game"
	"github.com/josephnaberhaus/gauthordle/internal/git"
	"github.com/josephnaberhaus/gauthordle/internal/output"
	"os"
	"strings"
)

const helpBody = "A daily game where you try to guess the author of some Git commits.\n\nTo play, simply \"git checkout\" the main development branch of your repository\nand run this program with no arguments.\n\nNew games start at midnight Central Time."

var (
	help   = flag.Bool("help", false, "Print the help message.")
	random = flag.Bool("random", false, "If true, play a random game instead of the daily game.")
)

func main() {
	flag.Parse()

	if len(flag.Args()) > 0 {
		exit(fmt.Errorf("unsupported arguments %q\n", strings.Join(flag.Args(), ",")))
	}

	if *help {
		showUsage()
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

func showUsage() {
	fmt.Println(helpBody)
	flag.Usage()

	os.Exit(0)
}

func exit(err error) {
	output.FprintColor(os.Stderr, fmt.Sprintf("ERROR: %s", err.Error()), output.Red)
	os.Exit(1)
}
