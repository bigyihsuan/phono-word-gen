package eval

import "phono-word-gen/parts"

type preparationResults struct {
	Categories map[string]parts.Category
}

func Prepare() (preparationResults, error) {
	results := preparationResults{
		Categories: make(map[string]parts.Category),
	}

	return results, nil
}
