package eval

import "math/rand"

type Category struct {
	Name     string
	Elements []CategoryElement
}

func NewCategory(name string, elements []CategoryElement) *Category {
	c := &Category{Name: name, Elements: elements}
	// TODO: set the weights of all phonemes
	// TODO: resolve embedded categories into weighted phonemes
	return c
}

func (c *Category) Get(random *rand.Rand) string {
	idx := random.Intn(len(c.Elements))
	return c.Elements[idx].Get(random)
}

type CategoryElement interface {
	Weight() float64
	SetWeight(float64)
	ResolveCategories(categories map[string]Category) []*Phoneme
	Get(random *rand.Rand) string
}

type UnresolvedCategory struct {
	Name     string
	Elements []CategoryElement
	weight   float64
}

func (uc UnresolvedCategory) Weight() float64           { return uc.weight }
func (uc *UnresolvedCategory) SetWeight(weight float64) { uc.weight = weight }
func (uc UnresolvedCategory) ResolveCategories(categories map[string]Category) []*Phoneme {
	//TODO
	return []*Phoneme{}
}
