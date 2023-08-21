package eval

import "honnef.co/go/js/dom/v2"

type elements struct {
	inputTextElement         *dom.HTMLTextAreaElement
	outputTextElement        *dom.HTMLTextAreaElement
	submitButton             *dom.HTMLButtonElement
	minSylCountElement       *dom.HTMLInputElement
	maxSylCountElement       *dom.HTMLInputElement
	wordCountElement         *dom.HTMLInputElement
	forbidDuplicatesElement  *dom.HTMLInputElement
	forceWordLimitElement    *dom.HTMLInputElement
	sortOutputElement        *dom.HTMLInputElement
	applyRejectionsElement   *dom.HTMLInputElement
	applyReplacementsElement *dom.HTMLInputElement
	generatedAlertElement    *dom.HTMLDivElement
	duplicateAlertElement    *dom.HTMLDivElement
	rejectedAlertElement     *dom.HTMLDivElement
	replacedAlertElement     *dom.HTMLDivElement
}

type entry struct {
	word Word
	syls []string
}
