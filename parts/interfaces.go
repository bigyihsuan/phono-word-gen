package parts

type Element interface {
	Get(categories map[string]Category) (string, error)
	ChoiceCount(categories map[string]Category) int
}

type SyllableElement interface {
	Element
	syllableElementTag()
}
