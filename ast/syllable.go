package ast

import (
	"fmt"
	"strings"
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

type SyllableSelection struct {
	Components []SyllableComponent
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
