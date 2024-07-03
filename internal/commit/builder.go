package commit

import (
	"cmp"
	"fmt"
	"github.com/josephnaberhaus/gauthordle/internal/config"
	"regexp"
	"time"
)

type FilterOption func(filter *Filter) error

func BuildFilter(options ...FilterOption) (*Filter, error) {
	filter := new(Filter)
	for _, opt := range options {
		err := opt(filter)
		if err != nil {
			return nil, fmt.Errorf("error build commit filter: %w", err)
		}
	}

	return filter, nil
}

func WithConfig(cfg config.Config) FilterOption {
	return func(filter *Filter) error {
		for _, authorFilter := range cfg.AuthorFilters {
			r, err := regexp.Compile(
				cmp.Or(authorFilter.ExcludeName, authorFilter.ExcludeEmail),
			)
			if err != nil {
				return err
			}

			if authorFilter.ExcludeName != "" {
				filter.nameFilters = append(filter.nameFilters, r)
			} else {
				filter.emailFilters = append(filter.emailFilters, r)
			}
		}

		return nil
	}
}

func WithStartTime(startTime time.Time) FilterOption {
	return func(filter *Filter) error {
		filter.startTime = startTime

		return nil
	}
}

func WithEndTime(endTime time.Time) FilterOption {
	return func(filter *Filter) error {
		filter.endTime = endTime

		return nil
	}
}
