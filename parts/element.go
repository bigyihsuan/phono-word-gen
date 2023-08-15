package parts

type Element interface {
	Get(categories map[string]Category) (string, error)
}
