package game

import (
	"github.com/JosephNaberhaus/gauthordle/internal/git"
	"github.com/JosephNaberhaus/prompt"
	"strconv"
	"time"
)

const header = "                   _   _                   _ _      \n  __ _  __ _ _   _| |_| |__   ___  _ __ __| | | ___ \n / _` |/ _` | | | | __| '_ \\ / _ \\| '__/ _` | |/ _ \\\n| (_| | (_| | |_| | |_| | | | (_) | | | (_| | |  __/\n \\__, |\\__,_|\\__,_|\\__|_| |_|\\___/|_|  \\__,_|_|\\___|\n |___/      The daily git author guessing game      "
const youWin = "                              _       \n _   _  ___  _   _  __      _(_)_ __  \n| | | |/ _ \\| | | | \\ \\ /\\ / / | '_ \\ \n| |_| | (_) | |_| |  \\ V  V /| | | | |\n \\__, |\\___/ \\__,_|   \\_/\\_/ |_|_| |_|\n |___/                                "
const youLose = "                     _                \n _   _  ___  _   _  | | ___  ___  ___ \n| | | |/ _ \\| | | | | |/ _ \\/ __|/ _ \\\n| |_| | (_) | |_| | | | (_) \\__ \\  __/\n \\__, |\\___/ \\__,_| |_|\\___/|___/\\___|\n |___/                                "
const nope = "                        \n _ __   ___  _ __   ___ \n| '_ \\ / _ \\| '_ \\ / _ \\\n| | | | (_) | |_) |  __/\n|_| |_|\\___/| .__/ \\___|\n            |_|         "

const numPuzzleCommits = 4

type puzzleHints struct {
	totalCommits    int
	mostTouchedFile string
}

type Puzzle struct {
	authorEmail   string
	authorName    string
	authorCommits []git.Commit
	puzzleCommits [numPuzzleCommits]git.Commit

	hints puzzleHints

	// All commits by all users.
	allCommits     []git.Commit
	allAuthorNames map[string]string
}

func (p Puzzle) Run() error {
	var promptOptions []prompt.SelectionOption
	for authorEmail, authorName := range p.allAuthorNames {
		promptOptions = append(promptOptions, prompt.SelectionOption{
			ID:          authorEmail,
			Name:        authorName,
			Description: authorEmail,
		})
	}

	for stage := 0; stage < numPuzzleCommits; stage++ {
		clearScreen()
		printColorLn(header, textColorYellow)
		ln()

		printColor("Guess the author of the following commit", textColorWhite)
		if stage > 0 {
			printColor("s", textColorYellow)
		}
		printColorLn(":", textColorYellow)
		ln()

		for i := 0; i <= stage; i++ {
			printColor("Commit #", textColorGreen)
			printColor(strconv.Itoa(i+1), textColorGreen)
			printColor(": ", textColorGreen)
			printColorLn(p.puzzleCommits[i].SubjectLine, textColorWhite)
		}

		// Hints
		ln()
		if stage >= 1 {
			ln()
			printColorLn("Hints", textColorGreen)
			printColor("Number of commits made by author in the last year: ", textColorGreen)
			printColorLn(strconv.Itoa(p.hints.totalCommits), textColorWhite)
		}
		if stage >= 3 {
			printColor("Author's most touched file: ", textColorGreen)
			printColorLn(p.hints.mostTouchedFile, textColorWhite)
		}

		ln()
		answerPrompt := &prompt.Select{
			Question:      "Who is the author?",
			Options:       promptOptions,
			NumLinesShown: 3,
		}
		err := answerPrompt.Show()
		if err != nil {
			return err
		}

		if answerPrompt.Response().ID == p.authorEmail {
			flashMessage(youWin, textColorGreen)
			break
		} else {
			if stage == numPuzzleCommits-1 {
				flashMessage(youLose, textColorRed)
			} else {
				flashMessage(nope, textColorRed)
			}
		}
	}

	ln()
	printColor("The answer was: ", textColorWhite)
	printColor(p.authorName, textColorWhite)
	printColor(" (", textColorWhite)
	printColor(p.authorEmail, textColorWhite)
	printColorLn(")", textColorWhite)
	ln()

	return nil
}

func flashMessage(message string, color textColor) {
	clearScreen()
	printColorLn(header, textColorYellow)
	ln()
	printColorLn(message, color)
	time.Sleep(1000 * time.Millisecond)
}
