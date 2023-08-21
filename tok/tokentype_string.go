// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package tok

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-1]
	_ = x[ILLEGAL-2]
	_ = x[COMMENT-3]
	_ = x[LPAREN-4]
	_ = x[RPAREN-5]
	_ = x[LBRACKET-6]
	_ = x[RBRACKET-7]
	_ = x[LBRACE-8]
	_ = x[RBRACE-9]
	_ = x[COMMA-10]
	_ = x[STAR-11]
	_ = x[COLON-12]
	_ = x[LINE_ENDING-13]
	_ = x[NUMBER-14]
	_ = x[EQ-15]
	_ = x[DOLLAR-16]
	_ = x[RAW-17]
	_ = x[CARET-18]
	_ = x[BSLASH-19]
	_ = x[AT-20]
	_ = x[AMPERSAND-21]
	_ = x[BANG-22]
	_ = x[PIPE-23]
	_ = x[ARROW-24]
	_ = x[SLASH-25]
	_ = x[DOUBLESLASH-26]
	_ = x[SYLLABLE-27]
	_ = x[LETTERS-28]
	_ = x[REJECT-29]
	_ = x[REPLACE-30]
}

const _TokenType_name = "EOFILLEGALCOMMENTLPARENRPARENLBRACKETRBRACKETLBRACERBRACECOMMASTARCOLONLINE_ENDINGNUMBEREQDOLLARRAWCARETBSLASHATAMPERSANDBANGPIPEARROWSLASHDOUBLESLASHSYLLABLELETTERSREJECTREPLACE"

var _TokenType_index = [...]uint8{0, 3, 10, 17, 23, 29, 37, 45, 51, 57, 62, 66, 71, 82, 88, 90, 96, 99, 104, 110, 112, 121, 125, 129, 134, 139, 150, 158, 165, 171, 178}

func (i TokenType) String() string {
	i -= 1
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
