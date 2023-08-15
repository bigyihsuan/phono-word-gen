package errs

import (
	"fmt"
	"phono-word-gen/tok"
)

var (
	ParserUnexpectedTokenError = func(got, want tok.TokenType) error { return fmt.Errorf("unexpected token: want=%q got=%q", want, got) }
)
