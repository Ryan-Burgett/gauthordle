package git

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterOutBots(t *testing.T) {
	input := []Commit{
		{
			AuthorName:  "Real person",
			AuthorEmail: "joe.smith@example.com",
		},
		{
			AuthorName:  "No reply",
			AuthorEmail: "joe.smith@example.noreply.com",
		},
		{
			AuthorName:  "NoSpace",
			AuthorEmail: "joe.smith@example.com",
		},
	}

	expected := []Commit{
		{
			AuthorName:  "Real person",
			AuthorEmail: "joe.smith@example.com",
		},
	}

	assert.Equal(t, expected, filterOutBots(input))
}

func TestConsolidateAuthorDetails(t *testing.T) {
	input := []Commit{
		{
			AuthorName:  "Joe Smith",
			AuthorEmail: "jOe.smith@example.com",
		},
		{
			AuthorName:  "Joseph Smith",
			AuthorEmail: "joe.smith@example.com",
		},
	}

	expected := []Commit{
		{
			AuthorName:  "Joseph Smith",
			AuthorEmail: "joe.smith@example.com",
		},
		{
			AuthorName:  "Joseph Smith",
			AuthorEmail: "joe.smith@example.com",
		},
	}

	assert.Equal(t, expected, consolidateAuthorDetails(input))
}

func TestFilterCommitSubjects(t *testing.T) {
	input := []Commit{
		{
			AuthorEmail: "joe.smith@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			AuthorEmail: "joe.smith@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			AuthorEmail: "joe.smith@example.com",
			SubjectLine: "DUPLICATED commit message",
		},
		{
			AuthorEmail: "bob.barker@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			SubjectLine: "oneword",
		},
		{
			SubjectLine: "two words",
		},
		{
			SubjectLine: "three words now",
		},
		{
			SubjectLine: "Merge branch 'master' of...",
		},
	}

	expected := []Commit{
		{
			AuthorEmail: "joe.smith@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			AuthorEmail: "bob.barker@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			SubjectLine: "three words now",
		},
	}

	assert.Equal(t, expected, filterCommitSubjects(input))
}
