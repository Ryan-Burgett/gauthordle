package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/josephnaberhaus/gauthordle/internal/git"
)

type builder struct {
	randomSource rand.Source
	commits      []git.Commit
	files        []string
	authorBias   float64
}

type Option func(*builder)

func WithRandomSource(randomSource rand.Source) Option {
	return func(b *builder) {
		b.randomSource = randomSource
	}
}

func WithCommits(commits []git.Commit) Option {
	return func(b *builder) {
		b.commits = commits
	}
}

func WithFiles(files []string) Option {
	return func(b *builder) {
		b.files = files
	}
}

func WithAuthorBias(authorBias float64) Option {
	return func(b *builder) {
		b.authorBias = authorBias
	}
}

// TODO add ability to generate game based on filetypes
func BuildPuzzle(opts ...Option) (Puzzle, error) {
	b := new(builder)
	for _, opt := range opts {
		opt(b)
	}

	// If no source of randomness was specified then just build a random one from the current timestamp.
	if b.randomSource == nil {
		b.randomSource = rand.NewSource(time.Now().Unix())
	}
	if b.authorBias < 1 || b.authorBias > 5 {
		return Puzzle{}, errors.New("author bias must be between 1 and 5")
	}

	return b.buildPuzzle()
}

func (b builder) buildPuzzle() (Puzzle, error) {
	random := rand.New(b.randomSource)

	author, err := pickAuthor(b.commits, b.authorBias, random)
	if err != nil {
		return Puzzle{}, fmt.Errorf("error building puzzle: %w", err)
	}

	authorNames := nameByEmail(b.commits)
	commitsByAuthor := commitsByAuthorEmail(b.commits)

	mostTouchedFile, err := mostTouchedFileForAuthor(author)
	if err != nil {
		return Puzzle{}, fmt.Errorf("error building puzzle: %w", err)
	}

	return Puzzle{
		authorEmail:   author,
		authorName:    authorNames[author],
		authorCommits: commitsByAuthor[author],
		puzzleCommits: pickPuzzleCommits(commitsByAuthor[author], random),
		hints: puzzleHints{
			totalCommits:    len(commitsByAuthor[author]),
			mostTouchedFile: mostTouchedFile,
		},
		allCommits:     b.commits,
		allAuthorNames: authorNames,
	}, nil
}

func pickPuzzleCommits(authorCommits []git.Commit, random *rand.Rand) [numPuzzleCommits]git.Commit {
	var result [numPuzzleCommits]git.Commit
	pickedIndices := map[int]struct{}{}
	for i := 0; i < numPuzzleCommits; i++ {
		index := random.Intn(len(authorCommits))
		if _, ok := pickedIndices[index]; ok {
			// We've already picked this number. Try again.
			i--
			continue
		}

		result[i] = authorCommits[index]
		pickedIndices[index] = struct{}{}
	}

	return result
}
