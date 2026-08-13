package ast

import (
	"fmt"
	"phono-word-gen/parts"
	"strings"

	"github.com/mroth/weightedrand/v3"
)

type SyllableDirective struct {
	Components []SyllableComponent
}

func (sd *SyllableDirective) node()      {}
func (sd *SyllableDirective) directive() {}
func (sd *SyllableDirective) String() string {
	components := []string{}
	for _, c := range sd.Components {
		components = append(components, c.String())
	}
	return fmt.Sprintf("(syllable %s)", strings.Join(components, " "))
}

type SyllableGrouping struct {
	Components []SyllableComponent
}

func (sg *SyllableGrouping) node()              {}
func (sg *SyllableGrouping) syllableComponent() {}
func (sg *SyllableGrouping) String() string {
	components := []string{}
	for _, c := range sg.Components {
		components = append(components, c.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(components, " "))
}

func (sg SyllableGrouping) SyllableElement() parts.SyllableElement {
	elements := []parts.SyllableElement{}
	for _, c := range sg.Components {
		elements = append(elements, c.SyllableElement())
	}
	return parts.NewGrouping(elements...)
}

type SyllableSelection struct {
	Components []WeightedSyllableComponent
}

func (ss *SyllableSelection) node()              {}
func (ss *SyllableSelection) syllableComponent() {}
func (ss *SyllableSelection) String() string {
	components := []string{}
	for _, c := range ss.Components {
		components = append(components, c.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(components, ", "))
}

func (ss SyllableSelection) SyllableElement() parts.SyllableElement {
	elements := []weightedrand.Choice[parts.SyllableElement, int]{}
	for _, c := range ss.Components {
		comp, weight := c.WeightedSyllableElement()
		choice := weightedrand.NewChoice(comp, weight)
		elements = append(elements, choice)
	}
	return parts.NewSelection(elements...)
}

type SyllableOptional struct {
	Components []SyllableComponent
	Weight     int
}

func (so *SyllableOptional) node()              {}
func (so *SyllableOptional) syllableComponent() {}
func (so *SyllableOptional) String() string {
	components := []string{}
	for _, c := range so.Components {
		components = append(components, c.String())
	}
	return fmt.Sprintf("((%s) * %d)", strings.Join(components, " "), so.Weight)
}

func (so SyllableOptional) SyllableElement() parts.SyllableElement {
	elements := []parts.SyllableElement{}
	for _, c := range so.Components {
		elements = append(elements, c.SyllableElement())
	}
	return parts.NewOptional(elements, so.Weight)
}

type WeightedSyllableComponent struct {
	Components []SyllableComponent
	Weight     int
}

func (wsc *WeightedSyllableComponent) node()              {}
func (wsc *WeightedSyllableComponent) syllableComponent() {}
func (wsc *WeightedSyllableComponent) String() string {
	components := []string{}
	for _, c := range wsc.Components {
		components = append(components, c.String())
	}
	return fmt.Sprintf("(%s * %d)", strings.Join(components, " "), wsc.Weight)
}

func (wsc WeightedSyllableComponent) SyllableElement() parts.SyllableElement {
	e, _ := wsc.WeightedSyllableElement()
	return e
}

func (wsc WeightedSyllableComponent) WeightedSyllableElement() (parts.SyllableElement, int) {
	elements := []parts.SyllableElement{}
	for _, c := range wsc.Components {
		elements = append(elements, c.SyllableElement())
	}
	return parts.NewGrouping(elements...), wsc.Weight
}
