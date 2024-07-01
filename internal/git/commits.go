package git

import (
	"errors"
	"fmt"
	"github.com/josephnaberhaus/gauthordle/internal/command"
	"strings"
	"time"
)

type Commit struct {
	AuthorName  string
	AuthorEmail string
	SubjectLine string
}

func GetCommits(start, end time.Time) ([]Commit, error) {
	const gitLogFormat = "%an\u001F%ae\u001F%s\u001E"
	result, err := command.Run("git", "log", "--since="+start.Format(time.DateOnly), "--until="+end.Format(time.DateOnly), "--format="+gitLogFormat)
	if err != nil {
		return nil, fmt.Errorf("error when getting git logs: %w", err)
	}

	// The rest of the code is written with the assumption that there are no new lines.
	result = strings.ReplaceAll(result, "\n", "")

	records := strings.Split(result, "\u001E")
	if len(records) == 0 {
		return nil, errors.New("unexpected response from git log")
	}
	// The last record will be an empty line.
	records = records[:len(records)-1]

	var commits []Commit
	for _, record := range records {
		fields := strings.Split(record, "\u001F")
		if len(fields) != 3 {
			return nil, errors.New("unexpected response from git log")
		}

		commits = append(commits, Commit{
			AuthorName:  fields[0],
			AuthorEmail: fields[1],
			SubjectLine: fields[2],
		})
	}

	commits = filterOutBots(commits)
	commits = consolidateAuthorDetails(commits)
	commits = filterCommitSubjects(commits)

	return commits, nil
}

// filterBots attempts to filter bot-made commits out.
// It's impossible to cover all cases here, but we'll make a best effort.
func filterOutBots(commits []Commit) []Commit {
	var result []Commit
	for _, commit := range commits {
		// E-mails with "noreply" in them are usually associated with bots.
		if strings.Contains(commit.AuthorEmail, "noreply") {
			continue
		}

		// Often robots won't have a proper "<First> <Last>" name (with a space separator).
		// This is going to have many false positives, but in the repos I tested it is more fun.
		if !strings.Contains(commit.AuthorName, " ") {
			continue
		}

		result = append(result, commit)
	}

	return result
}

// consolidateAuthorNames ensures that each e-mail address is associated with only one author name.
func consolidateAuthorDetails(commits []Commit) []Commit {
	result := make([]Commit, len(commits))
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

// filterCommitSubjects ensures that we only use interesting commit subjects
func filterCommitSubjects(commits []Commit) []Commit {
	var result []Commit
	authorCommitSubjects := map[string]map[string]struct{}{}
	for _, commit := range commits {
		// Commit messages with only 1-2 words don't get you very much information.
		if len(strings.Split(commit.SubjectLine, " ")) < 3 {
			continue
		}

		// Some authors use the same commit message over and over again.
		// Remove the duplicates
		if authorCommitSubjects[commit.AuthorEmail] == nil {
			authorCommitSubjects[commit.AuthorEmail] = map[string]struct{}{}
		}

		// Merge commits aren't helpful
		if strings.HasPrefix(commit.SubjectLine, "Merge branch") {
			continue
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

// GetFilesChangedForAuthor gets all the files touched be the given user.
// The returned list can contain duplicates.
func GetFilesChangedForAuthor(authorEmail string) ([]string, error) {
	result, err := command.Run("git", "log", "--author="+authorEmail, "--name-only", "--format=")
	if err != nil {
		return nil, fmt.Errorf("error when getting files : %w", err)
	}

	return strings.Split(result, "\n"), nil
}
