// Code generated by "stringer -type=ContextPrefix"; DO NOT EDIT.

package parts

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NO_PREFIX-0]
	_ = x[WORD_START-1]
	_ = x[SYL_START-2]
	_ = x[NOT-3]
}

const _ContextPrefix_name = "NO_PREFIXWORD_STARTSYL_STARTNOT"

var _ContextPrefix_index = [...]uint8{0, 9, 19, 28, 31}

func (i ContextPrefix) String() string {
	if i < 0 || i >= ContextPrefix(len(_ContextPrefix_index)-1) {
		return "ContextPrefix(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ContextPrefix_name[_ContextPrefix_index[i]:_ContextPrefix_index[i+1]]
}
