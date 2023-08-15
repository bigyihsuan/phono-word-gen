package parts

import (
	"errors"
	"math/rand"
	"phono-word-gen/errs"
	"strings"

	wr "github.com/mroth/weightedrand/v2"
)

type Syllable struct {
	Elements []SyllableElement
}

func NewSyllable(elements ...SyllableElement) *Syllable { return &Syllable{Elements: elements} }
func (s *Syllable) Get(categories map[string]Category) string {
	elements := []string{}
	for _, e := range s.Elements {
		elements = append(elements, e.Get(categories))
	}
	return strings.Join(elements, "")
}

type SyllableElement interface {
	Element
	syllableElementTag()
}

// A raw string in a syllable
type Raw struct {
	Value string
}

func NewRaw(value string) *Raw                  { return &Raw{Value: value} }
func (r *Raw) syllableElementTag()              {}
func (r *Raw) Get(_ map[string]Category) string { return r.Value }

type Grouping struct {
	Elements []SyllableElement
}

func NewGrouping(elements ...SyllableElement) *Grouping { return &Grouping{Elements: elements} }
func (g *Grouping) syllableElementTag()                 {}
func (g *Grouping) Get(categories map[string]Category) string {
	// evaluate all elements in the grouping
	values := []string{}
	for _, v := range g.Elements {
		values = append(values, v.Get(categories))
	}
	return strings.Join(values, "")
}

type Selection struct {
	Choices *wr.Chooser[SyllableElement, int]
}

func NewSelection(elements ...wr.Choice[SyllableElement, int]) (*Selection, error) {
	chooser, err := wr.NewChooser(elements...)
	if err != nil {
		return nil, errors.Join(errs.SelectionCreationError, err)
	}
	return &Selection{Choices: chooser}, nil
}
func (s *Selection) syllableElementTag() {}
func (s *Selection) Get(catgories map[string]Category) string {
	// pick a random choice in the selection
	return s.Choices.Pick().Get(catgories)
}

type Optional struct {
	Elements []SyllableElement
	Chance   int // defaults to 50, for 50%
}

func NewOptional(elements []SyllableElement, chances ...int) *Optional {
	chance := 50
	if len(chances) > 0 {
		chance = chances[0]
	}
	return &Optional{Elements: elements, Chance: chance}
}
func (o *Optional) syllableElementTag() {}
func (o *Optional) Get(categories map[string]Category) string {
	// evaluate all elements in the optional if the rng is less than the chance
	if rand.Intn(101) < o.Chance {
		values := []string{}
		for _, v := range o.Elements {
			values = append(values, v.Get(categories))
		}
		return strings.Join(values, "")
	} else {
		return ""
	}
}
