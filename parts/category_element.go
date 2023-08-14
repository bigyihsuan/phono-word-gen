package parts

type CategoryElement interface {
	Get(categories map[string]Category) string
}
