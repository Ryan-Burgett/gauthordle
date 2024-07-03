package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/josephnaberhaus/gauthordle/internal/git"
)

// Increment to permanently change the RNG of the game.
const seed = 0

func BuildRandom() (Puzzle, error) {
	return buildGame(rand.Intn)
}

func BuildToday() (Puzzle, error) {
	startTime, _ := puzzleTimeRange()
	// Make a random based on the end time so that it's stable throughout the day.
	random := rand.New(rand.NewSource(startTime.Unix() + seed))
	return buildGame(random.Intn)
}

func buildGame(r func(int) int) (Puzzle, error) {
	startTime, endTime := puzzleTimeRange()
	commits, err := git.GetCommits(startTime, endTime)
	if err != nil {
		return Puzzle{}, fmt.Errorf("error building puzzle: %w", err)
	}

	author, err := pickAuthor(commits, r)
	if err != nil {
		return Puzzle{}, fmt.Errorf("error building puzzle: %w", err)
	}

	authorNames := nameByEmail(commits)
	commitsByAuthor := commitsByAuthorEmail(commits)

	mostTouchedFile, err := mostTouchedFileForAuthor(author)
	if err != nil {
		return Puzzle{}, fmt.Errorf("error building puzzle: %w", err)
	}

	return Puzzle{
		authorEmail:   author,
		authorName:    authorNames[author],
		authorCommits: commitsByAuthor[author],
		puzzleCommits: pickPuzzleCommits(commitsByAuthor[author], r),
		hints: puzzleHints{
			totalCommits:    len(commitsByAuthor[author]),
			mostTouchedFile: mostTouchedFile,
		},
		allCommits:     commits,
		allAuthorNames: authorNames,
	}, nil
}

func pickPuzzleCommits(authorCommits []git.Commit, random func(int) int) [numPuzzleCommits]git.Commit {
	var result [numPuzzleCommits]git.Commit
	pickedIndices := map[int]struct{}{}
	for i := 0; i < numPuzzleCommits; i++ {
		index := random(len(authorCommits))
		if _, ok := pickedIndices[index]; ok {
			// We've already picked this number. Try again.
			i--
			continue
		}

		result[i] = authorCommits[i]
		pickedIndices[i] = struct{}{}
	}

	return result
}

func puzzleTimeRange() (time.Time, time.Time) {
	// Center the games around Central Time (because that's where I live).
	gmtNow := time.Now().In(time.FixedZone("CT", 0))

	// End with commits from a year ago.
	startDate := time.Date(gmtNow.Year()-1, gmtNow.Month(), gmtNow.Day(), 0, 0, 0, 0, gmtNow.Location())
	// End with commits from a week ago to increase the odds that our user will have an up-to-date history.
	endDate := time.Date(gmtNow.Year(), gmtNow.Month(), gmtNow.Day()-7, 0, 0, 0, 0, gmtNow.Location())

	return startDate, endDate
}
