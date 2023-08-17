package parts

import (
	"errors"
	"fmt"
	"phono-word-gen/tok"
)

var (
	UnexpectedToken = func(token tok.Token, tts ...tok.TokenType) error {
		switch {
		case len(tts) < 2:
			return fmt.Errorf("unexpected token: got=%q (%s)", tts[0], token.String())
		case len(tts) == 2:
			return fmt.Errorf("unexpected token: got=%q (%s) want=%q", tts[0], token.String(), tts[1])
		default:
			return fmt.Errorf("unexpected token: got=%q (%s) want=%q", tts[0], token.String(), tts[1])
		}
	}
	UnknownDirective = func(tt tok.TokenType) error { return fmt.Errorf("unknown directive: got=%q", tt) }
)

var (
	RecursiveCategoryError = func(cat, ref string) error {
		return fmt.Errorf("recursive category: %s contains %s contains %s", cat, ref, cat)
	}
	UndefinedCategoryError = func(cat, ref string) error {
		return fmt.Errorf("undefined category: %s (contained in %s)", ref, cat)
	}
	CategoryCreationError  = errors.New("category creation error")
	SelectionCreationError = errors.New("selection creation error")
	OptionalCreationError  = errors.New("optional creation error")
)
