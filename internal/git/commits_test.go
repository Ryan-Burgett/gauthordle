package git

import (
	"slices"
	"testing"

	"github.com/josephnaberhaus/gauthordle/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestFilterForConfig(t *testing.T) {
	input := []Commit{
		{
			AuthorName:  "Person 1",
			AuthorEmail: "abc@abc.com",
		},
		{
			AuthorName:  "Person 2",
			AuthorEmail: "def@def.com",
		},
		{
			AuthorName:  "Person 3",
			AuthorEmail: "ghi@ghi.com",
		},
	}

	tests := []struct {
		desc string
		cfg  config.Config
		exp  []Commit
	}{{
		desc: "empty config should keep all",
		cfg:  config.Config{},
		exp:  input,
	}, {
		desc: "filter by name",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeName: "2",
			}},
		},
		exp: []Commit{input[0], input[2]},
	}, {
		desc: "filter by email",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeEmail: "ghi",
			}},
		},
		exp: []Commit{input[0], input[1]},
	}, {
		desc: "specify name and email, prefer name",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeName:  "2",
				ExcludeEmail: "ghi",
			}},
		},
		exp: []Commit{input[0], input[2]},
	}, {
		desc: "specify multiple filters",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeName: "2",
			}, {
				ExcludeEmail: "ghi",
			}},
		},
		exp: []Commit{input[0]},
	}, {
		desc: "filter by name regexp",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeName: "[2-3]",
			}},
		},
		exp: []Commit{input[0]},
	}, {
		desc: "filter by email regexp",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeEmail: "(abc|def)",
			}},
		},
		exp: []Commit{input[2]},
	}}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			commits, err := filterForConfig(slices.Clone(input), tc.cfg)
			require.NoError(t, err)
			assert.Equal(t, tc.exp, commits)
		})
	}
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
		{
			SubjectLine: "PREFIX Merge branch 'master' of...",
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
