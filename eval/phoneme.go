package eval

import (
	"math"
	"math/rand"
)

type Phoneme struct {
	Value  string
	weight float64 // NaN is the "unweighted" weight
}

func NewPhoneme(value string, weights ...float64) *Phoneme {
	p := Phoneme{Value: value, weight: math.NaN()}
	if len(weights) > 0 {
		p.weight = weights[0]
	}
	return &p
}

func (p *Phoneme) ResolveCategories(categories map[string]Category) []*Phoneme { return []*Phoneme{p} }
func (p *Phoneme) Weight() float64                                             { return p.weight }
func (p *Phoneme) SetWeight(weight float64)                                    { p.weight = weight }
func (p *Phoneme) Get(_ *rand.Rand) string                                     { return p.Value }
