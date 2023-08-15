package eval

import (
	"phono-word-gen/parts"

	"honnef.co/go/js/dom/v2"
)

type Evaluator struct {
	document dom.Document

	minSylCount, maxSylCount int
	wordCount                int
}

type preparationResults struct {
	Categories map[string]parts.Category
}

func Prepare() (preparationResults, error) {
	results := preparationResults{
		Categories: make(map[string]parts.Category),
	}

	return results, nil
}
