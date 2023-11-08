package parts

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"

	wr "github.com/mroth/weightedrand/v2"
)

type Syllable struct {
	Elements SyllableElements
}

func NewSyllable(elements ...SyllableElement) *Syllable { return &Syllable{Elements: elements} }
func (s *Syllable) Get(categories Categories, components Components) (string, error) {
	elements := []string{}
	for _, e := range s.Elements {
		ele, err := e.Get(categories, components)
		if err != nil {
			return ele, err
		}
		elements = append(elements, ele)
	}
	return strings.Join(elements, ""), nil
}
func (s *Syllable) ChoiceCount(categories Categories, components Components) int {
	count := 1
	for _, e := range s.Elements {
		count *= e.ChoiceCount(categories, components)
	}
	return count
}

type Grouping struct {
	Elements SyllableElements
}

func NewGrouping(elements ...SyllableElement) *Grouping { return &Grouping{Elements: elements} }
func (g *Grouping) syllableElementTag()                 {}
func (g *Grouping) Get(categories Categories, components Components) (string, error) {
	// evaluate all elements in the grouping
	values := []string{}
	for _, v := range g.Elements {
		val, err := v.Get(categories, components)
		if err != nil {
			return val, err
		}
		values = append(values, val)
	}
	return strings.Join(values, ""), nil
}
func (g *Grouping) ChoiceCount(categories Categories, components Components) int {
	count := 1
	for _, e := range g.Elements {
		count *= e.ChoiceCount(categories, components)
	}
	return count
}
func (g *Grouping) Regexp(categories Categories, components Components) *regexp.Regexp {
	elements := []string{}
	for _, e := range g.Elements {
		elements = append(elements, e.Regexp(categories, components).String())
	}
	return regexp.MustCompile("(" + strings.Join(elements, "") + ")")
}

type Selection struct {
	Choices []wr.Choice[SyllableElement, int]
}

func NewSelection(elements ...wr.Choice[SyllableElement, int]) *Selection {
	return &Selection{Choices: elements}
}
func (s *Selection) syllableElementTag() {}
func (s *Selection) Get(catgories Categories, components Components) (string, error) {
	// pick a random choice in the selection
	chooser, err := wr.NewChooser(s.Choices...)
	if err != nil {
		return "", errors.Join(SelectionCreationError, err)
	}
	return chooser.Pick().Get(catgories, components)
}
func (s *Selection) ChoiceCount(categories Categories, components Components) int {
	count := len(s.Choices)
	for _, choice := range s.Choices {
		count *= choice.Item.ChoiceCount(categories, components)
	}
	return count
}
func (s *Selection) Regexp(categories Categories, components Components) *regexp.Regexp {
	elements := []string{}
	for _, e := range s.Choices {
		elements = append(elements, e.Item.Regexp(categories, components).String())
	}
	return regexp.MustCompile("(" + strings.Join(elements, "|") + ")")
}

// optional component. defaults to 50% chance of appearing when calling Get().
type Optional struct {
	Elements SyllableElements
	weight   int
}

func NewOptional(elements SyllableElements, percentChance ...int) *Optional {
	weight := 50 // default to 50/50
	if len(percentChance) > 0 {
		weight = percentChance[0]
	}
	return &Optional{Elements: elements, weight: weight}
}
func (o *Optional) syllableElementTag() {}
func (o *Optional) Get(categories Categories, components Components) (string, error) {
	chance := rand.Intn(101)
	if chance < o.weight {
		return NewGrouping(o.Elements...).Get(categories, components)
	} else {
		return "", nil
	}
}
func (o *Optional) ChoiceCount(categories Categories, components Components) int {
	count := 2
	for _, e := range o.Elements {
		count *= e.ChoiceCount(categories, components)
	}
	return count
}
func (o *Optional) Regexp(categories Categories, components Components) *regexp.Regexp {
	elements := []string{}
	for _, e := range o.Elements {
		elements = append(elements, e.Regexp(categories, components).String())
	}
	return regexp.MustCompile("(" + strings.Join(elements, "") + ")?")
}
