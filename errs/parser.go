package errs

import (
	"fmt"
	"phono-word-gen/tok"
)

var (
	UnexpectedToken = func(tts ...tok.TokenType) error {
		switch {
		case len(tts) < 2:
			return fmt.Errorf("unexpected token: got=%q", tts[0])
		case len(tts) == 2:
			return fmt.Errorf("unexpected token: got=%q want=%q", tts[0], tts[1])
		default:
			return fmt.Errorf("unexpected token: got=%q want=%q", tts[0], tts[1])
		}
	}
	UnknownDirective = func(tt tok.TokenType) error { return fmt.Errorf("unknown directive: got=%q", tt) }
)
