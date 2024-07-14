package file

import (
	"cmp"
	"github.com/josephnaberhaus/gauthordle/internal/config"
	"regexp"
)

type FilterOption func(filter *Filter) error

func WithTeam(s string) FilterOption {
	//TODO implement this
	return nil
}

func BuildFilter(options ...FilterOption) (Filter, error) {
	return Filter{}, nil
}

func WithConfig(cfg config.Config) FilterOption {
	return func(filter *Filter) error {
		for _, fileTypeFilter := range cfg.FileTypeFilters {
			r, err := regexp.Compile(
				cmp.Or(fileTypeFilter.ExcludeFileType),
			)
			if err != nil {
				return err
			}

			if fileTypeFilter.ExcludeFileType != "" {
				filter.typeFilters = append(filter.typeFilters, r)
			}
		}
		return nil
	}
}
