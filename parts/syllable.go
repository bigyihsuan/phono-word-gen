package parts

import (
	"errors"
	"phono-word-gen/errs"
	"strings"

	wr "github.com/mroth/weightedrand/v2"
)

type Syllable struct {
	Elements []SyllableElement
}

func NewSyllable(elements ...SyllableElement) *Syllable { return &Syllable{Elements: elements} }
func (s *Syllable) Get(categories map[string]Category) (string, error) {
	elements := []string{}
	for _, e := range s.Elements {
		ele, err := e.Get(categories)
		if err != nil {
			return ele, err
		}
		elements = append(elements, ele)
	}
	return strings.Join(elements, ""), nil
}

type SyllableElement interface {
	Element
	syllableElementTag()
}

// A raw string in a syllable
type Raw struct {
	Value string
}

func NewRaw(value string) *Raw                           { return &Raw{Value: value} }
func (r *Raw) syllableElementTag()                       {}
func (r *Raw) Get(_ map[string]Category) (string, error) { return r.Value, nil }

type Grouping struct {
	Elements []SyllableElement
}

func NewGrouping(elements ...SyllableElement) *Grouping { return &Grouping{Elements: elements} }
func (g *Grouping) syllableElementTag()                 {}
func (g *Grouping) Get(categories map[string]Category) (string, error) {
	// evaluate all elements in the grouping
	values := []string{}
	for _, v := range g.Elements {
		val, err := v.Get(categories)
		if err != nil {
			return val, err
		}
		values = append(values, val)
	}
	return strings.Join(values, ""), nil
}

type Selection struct {
	Choices []wr.Choice[SyllableElement, int]
}

func NewSelection(elements ...wr.Choice[SyllableElement, int]) *Selection {
	return &Selection{Choices: elements}
}
func (s *Selection) syllableElementTag() {}
func (s *Selection) Get(catgories map[string]Category) (string, error) {
	// pick a random choice in the selection
	chooser, err := wr.NewChooser(s.Choices...)
	if err != nil {
		return "", errors.Join(errs.SelectionCreationError, err)
	}
	return chooser.Pick().Get(catgories)
}

// optional component. defaults to 50% chance of appearing when calling Get().
type Optional struct {
	Elements []SyllableElement
	weight   int
}

func NewOptional(elements []SyllableElement, percentChance ...int) *Optional {
	weight := 50 // default to 50/50
	if len(percentChance) > 0 {
		weight = percentChance[0]
	}
	return &Optional{Elements: elements, weight: weight}
}
func (o *Optional) syllableElementTag() {}
func (o *Optional) Get(categories map[string]Category) (string, error) {
	chooser, err := wr.NewChooser[SyllableElement, int](
		wr.NewChoice[SyllableElement, int](NewGrouping(o.Elements...), 100-o.weight),
		wr.NewChoice[SyllableElement, int](nil, o.weight),
	)
	if err != nil {
		return "", errors.Join(errs.OptionalCreationError, err)
	}
	element := chooser.Pick()
	if element == nil {
		return "", nil
	} else {
		return element.Get(categories)
	}
}
