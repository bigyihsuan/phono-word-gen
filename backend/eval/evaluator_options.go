package eval

import "syscall/js"

type Options struct {
	Phonology     string `json:"phonology"`
	MinSylCount   int    `json:"minSylCount"`
	MaxSylCount   int    `json:"maxSylCount"`
	WordCount     int    `json:"wordCount"`
	SentenceCount int    `json:"sentenceCount"`

	ForbidDuplicates  bool `json:"forbidDuplicates"`
	ForceWordLimit    bool `json:"forceWordLimit"`
	SortOutput        bool `json:"sortOutput"`
	MarkSyllables     bool `json:"markSyllables"`
	ApplyRejections   bool `json:"applyRejections"`
	ApplyReplacements bool `json:"applyReplacements"`
	GenerateSentences bool `json:"generateSentences"`
}

const GENERATE_WORDS = "words"
const GENERATE_SENTENCES = "sentences"

func OptionsFromJsValue(v js.Value) Options {
	o := Options{}
	o.Phonology = v.Get("phonology").String()
	o.MinSylCount = v.Get("minSylCount").Int()
	o.MaxSylCount = v.Get("maxSylCount").Int()
	o.WordCount = v.Get("wordCount").Int()
	o.SentenceCount = v.Get("sentenceCount").Int()
	o.ForbidDuplicates = v.Get("forbidDuplicates").Bool()
	o.ForceWordLimit = v.Get("forceWordLimit").Bool()
	o.SortOutput = v.Get("sortOutput").Bool()
	o.MarkSyllables = v.Get("markSyllables").Bool()
	o.ApplyRejections = v.Get("applyRejections").Bool()
	o.ApplyReplacements = v.Get("applyReplacements").Bool()
	o.GenerateSentences = v.Get("generateSentences").String() == GENERATE_WORDS
	return o
}
