package parts

import (
	"errors"
	"math/rand"
	"regexp"
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
func (s *Syllable) ChoiceCount(categories map[string]Category) int {
	count := 1
	for _, e := range s.Elements {
		count *= e.ChoiceCount(categories)
	}
	return count
}

// A raw string in a syllable
type Raw struct {
	Value string
}

func NewRaw(value string) *Raw                                { return &Raw{Value: value} }
func (r *Raw) syllableElementTag()                            {}
func (r *Raw) Get(_ map[string]Category) (string, error)      { return r.Value, nil }
func (r *Raw) ChoiceCount(categories map[string]Category) int { return 1 }
func (r *Raw) Regexp(_ map[string]Category) *regexp.Regexp {
	return regexp.MustCompile("(" + regexp.QuoteMeta(r.Value) + ")")
}

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
func (g *Grouping) ChoiceCount(categories map[string]Category) int {
	count := 1
	for _, e := range g.Elements {
		count *= e.ChoiceCount(categories)
	}
	return count
}
func (g *Grouping) Regexp(categories map[string]Category) *regexp.Regexp {
	elements := []string{}
	for _, e := range g.Elements {
		elements = append(elements, e.Regexp(categories).String())
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
func (s *Selection) Get(catgories map[string]Category) (string, error) {
	// pick a random choice in the selection
	chooser, err := wr.NewChooser(s.Choices...)
	if err != nil {
		return "", errors.Join(SelectionCreationError, err)
	}
	return chooser.Pick().Get(catgories)
}
func (s *Selection) ChoiceCount(categories map[string]Category) int {
	count := len(s.Choices)
	for _, choice := range s.Choices {
		count *= choice.Item.ChoiceCount(categories)
	}
	return count
}
func (s *Selection) Regexp(categories map[string]Category) *regexp.Regexp {
	elements := []string{}
	for _, e := range s.Choices {
		elements = append(elements, e.Item.Regexp(categories).String())
	}
	return regexp.MustCompile("(" + strings.Join(elements, "|") + ")")
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
	chance := rand.Intn(101)
	if chance < o.weight {
		return NewGrouping(o.Elements...).Get(categories)
	} else {
		return "", nil
	}
}
func (o *Optional) ChoiceCount(categories map[string]Category) int {
	count := 2
	for _, e := range o.Elements {
		count *= e.ChoiceCount(categories)
	}
	return count
}
func (o *Optional) Regexp(categories map[string]Category) *regexp.Regexp {
	elements := []string{}
	for _, e := range o.Elements {
		elements = append(elements, e.Regexp(categories).String())
	}
	return regexp.MustCompile("(" + strings.Join(elements, "") + ")?")
}
