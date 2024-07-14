package file

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
)

type Filter struct {
	// typeFilters specifies what file types should be excluded.
	typeFilters []*regexp.Regexp
}

func (f *Filter) GetFiles() ([]string, error) {
	files, err := getFiles()
	if err != nil {
		return nil, err
	}

	files = f.filterExclusions(files)
	return files, nil
}

func getFiles() ([]string, error) {
	files := make([]string, 0)
	// walks the file path from the current directory and subdirectories
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {

				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// filterExclusions filters out excluded file types
func (f *Filter) filterExclusions(files []string) []string {
	shouldDelete := func(file string) bool {
		for _, f := range f.typeFilters {
			if f.MatchString(path.Ext(file)) {
				return true
			}
		}
		return false
	}

	files = slices.Clone(files)
	return slices.DeleteFunc(files, shouldDelete)
}
