package git

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/josephnaberhaus/gauthordle/internal/command"
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

	return commits, nil
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
