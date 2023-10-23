package parts

//go:generate stringer -type=ContextPrefix
//go:generate stringer -type=ContextSuffix
type ContextPrefix int
type ContextSuffix int

const (
	NO_PREFIX ContextPrefix = iota
	WORD_START
	SYL_START
	NOT
)
const (
	NO_SUFFIX ContextSuffix = iota
	WORD_END
	SYL_END
)
