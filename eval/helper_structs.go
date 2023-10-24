package eval

import "honnef.co/go/js/dom/v2"

type Elements struct {
	inputTextElement         *dom.HTMLTextAreaElement
	outputTextElement        *dom.HTMLTextAreaElement
	submitButtonElement      *dom.HTMLButtonElement
	minSylCountElement       *dom.HTMLInputElement
	maxSylCountElement       *dom.HTMLInputElement
	wordCountElement         *dom.HTMLInputElement
	sentenceCountElement     *dom.HTMLInputElement
	generateSentencesElement *dom.HTMLInputElement
	forbidDuplicatesElement  *dom.HTMLInputElement
	forceWordLimitElement    *dom.HTMLInputElement
	sortOutputElement        *dom.HTMLInputElement
	markSyllablesElement     *dom.HTMLInputElement
	applyRejectionsElement   *dom.HTMLInputElement
	applyReplacementsElement *dom.HTMLInputElement
	generatedAlertElement    *dom.HTMLDivElement
	duplicateAlertElement    *dom.HTMLDivElement
	rejectedAlertElement     *dom.HTMLDivElement
	replacedAlertElement     *dom.HTMLDivElement
	copyButtonElement        *dom.HTMLButtonElement
}

type Options struct {
	minSylCount, maxSylCount int
	wordCount, sentenceCount int

	forbidDuplicates, forceWordLimit, sortOutput, markSyllables bool
	applyRejections, applyReplacements, generateSentences       bool
}

type entry struct {
	word Word
	syls []string
}
