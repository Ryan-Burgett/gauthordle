package commit

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/josephnaberhaus/gauthordle/internal/git"
)

type Filter struct {
	// startTime specifies the oldest commit to return.
	// endTime specifies the earliest commit to return.
	startTime, endTime time.Time
	// nameFilters specifies what author names should be excluded.
	nameFilters []*regexp.Regexp
	// emailFilters specifies what author emails should be excluded.
	emailFilters []*regexp.Regexp
	// team specifies what team to return commits for.
	// If this doesn't match an element of teams than all commits will be returned.
	team string
	// teams is a map from team name to a set of e-mails for the members of the team.
	teams map[string]map[string]struct{}
}

func (f *Filter) GetCommits() ([]git.Commit, error) {
	if f.startTime.IsZero() {
		return nil, fmt.Errorf("no start time specified")
	}
	if f.endTime.IsZero() {
		return nil, fmt.Errorf("no end time specified")
	}

	commits, err := git.GetCommits(f.startTime, f.endTime)
	if err != nil {
		return nil, err
	}

	type filterFunc func([]git.Commit) []git.Commit
	filters := []filterFunc{
		f.filterExclusions,
		f.filterByTeam,
		f.filterOutBots,
		f.filterCommitSubjects,
		f.consolidateAuthorDetails,
	}
	for _, filter := range filters {
		commits = filter(commits)
	}

	return commits, nil
}

// filterBots attempts to filter bot-made commits out.
// It's impossible to cover all cases here, but we'll make a best effort.
func (f *Filter) filterOutBots(commits []git.Commit) []git.Commit {
	var result []git.Commit
	for _, commit := range commits {
		// E-mails with "noreply" in them are usually associated with bots.
		if strings.Contains(commit.AuthorEmail, "noreply") {
			continue
		}

		// If the author's name contains the word "robot" that's a pretty good indication that it's a robot.
		if strings.Contains(strings.ToLower(commit.AuthorName), "robot") {
			continue
		}

		result = append(result, commit)
	}

	return result
}

// filterExclusions filters out excluded names and e-mails.
func (f *Filter) filterExclusions(commits []git.Commit) []git.Commit {
	shouldDelete := func(commit git.Commit) bool {
		for _, f := range f.nameFilters {
			if f.MatchString(commit.AuthorName) {
				return true
			}
		}
		for _, f := range f.emailFilters {
			if f.MatchString(commit.AuthorEmail) {
				return true
			}
		}
		return false
	}

	commits = slices.Clone(commits)
	return slices.DeleteFunc(commits, shouldDelete)
}

func (f *Filter) filterByTeam(commits []git.Commit) []git.Commit {
	team, ok := f.teams[f.team]
	if !ok {
		// If not valid team is specified then just use all commits.
		return commits
	}

	var result []git.Commit
	for _, commit := range commits {
		if _, ok := team[strings.ToLower(commit.AuthorEmail)]; ok {
			result = append(result, commit)
		}
	}

	return result
}

// filterCommitSubjects ensures that we only use interesting commit subjects
func (f *Filter) filterCommitSubjects(commits []git.Commit) []git.Commit {
	var result []git.Commit
	authorCommitSubjects := map[string]map[string]struct{}{}
	for _, commit := range commits {
		// Commit messages with only 1-2 words don't get you very much information.
		if len(strings.Split(commit.SubjectLine, " ")) < 3 {
			continue
		}

		// Merge commits aren't helpful
		if strings.Contains(strings.ToLower(commit.SubjectLine), "merge") {
			continue
		}

		// Some authors use the same commit message over and over again.
		// Remove the duplicates
		if authorCommitSubjects[commit.AuthorEmail] == nil {
			authorCommitSubjects[commit.AuthorEmail] = map[string]struct{}{}
		}

		// No reason for this to be case sensitive
		subjectLine := strings.ToLower(commit.SubjectLine)

		authorCommits := authorCommitSubjects[commit.AuthorEmail]
		if _, ok := authorCommits[subjectLine]; ok {
			continue
		}
		authorCommits[subjectLine] = struct{}{}

		result = append(result, commit)
	}

	return result
}

// consolidateAuthorDetails ensures that the e-mails and names for all authors are uniform.
func (f *Filter) consolidateAuthorDetails(commits []git.Commit) []git.Commit {
	result := make([]git.Commit, len(commits))
	emailToAuthorName := map[string]string{}

	// Go backwards so that we use the most recently used name
	for i := len(commits) - 1; i >= 0; i-- {
		commit := commits[i]

		// E-mails are not case-sensitive, so lower-case it.
		authorEmail := strings.ToLower(commit.AuthorEmail)
		commit.AuthorEmail = authorEmail

		if _, ok := emailToAuthorName[authorEmail]; !ok {
			// This is the first time we have encountered this e-mail.
			emailToAuthorName[authorEmail] = commit.AuthorName
		} else {
			// We've encountered this e-mail before. Make sure we use that same name.
			commit.AuthorName = emailToAuthorName[authorEmail]
		}

		result[i] = commit
	}

	return result
}
