package errs

import (
	"fmt"
	"phono-word-gen/tok"
)

var (
	ParserUnexpectedToken = func(tts ...tok.TokenType) error {
		switch {
		case len(tts) < 2:
			return fmt.Errorf("unexpected token: got=%q", tts[0])
		case len(tts) == 2:
			return fmt.Errorf("unexpected token: got=%q want=%q", tts[0], tts[1])
		default:
			return fmt.Errorf("unexpected token: got=%q want=%q", tts[0], tts[1])
		}
	}
	ParserUnknownDirective = func(tt tok.TokenType) error { return fmt.Errorf("unknwown directive: got=%q", tt) }
)
