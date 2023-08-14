package eval

type CategoryElement interface {
	Get(categories map[string]Category) string
}
