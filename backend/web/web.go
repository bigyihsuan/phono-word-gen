package web

import (
	"errors"
	"fmt"
	"io"
	"phono-word-gen/eval"
	"phono-word-gen/sample"
	"phono-word-gen/util"
	"strconv"
	"strings"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

type Web struct {
	document dom.Document
	Elements
	Evaluator *eval.Evaluator

	*examplePageElements
}

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

type entry struct {
	word eval.Word
	syls []string
}

type examplePageElements struct {
	sampleDropdownElement *dom.HTMLSelectElement
}

func New() (*Web, error) {
	w := new(Web)
	w.loadDocument()
	w.setEventListeners()
	w.Evaluator = eval.New(w.getOptions())
	return w, nil
}

func (w *Web) loadDocument() {
	w.document = dom.GetWindow().Document()
	w.inputTextElement = w.document.QuerySelector("#phonology").(*dom.HTMLTextAreaElement)
	w.outputTextElement = w.document.QuerySelector("#outputText").(*dom.HTMLTextAreaElement)
	w.submitButtonElement = w.document.QuerySelector("#submit").(*dom.HTMLButtonElement)
	w.minSylCountElement = w.document.QuerySelector("#minSylCount").(*dom.HTMLInputElement)
	w.maxSylCountElement = w.document.QuerySelector("#maxSylCount").(*dom.HTMLInputElement)
	w.wordCountElement = w.document.QuerySelector("#wordCount").(*dom.HTMLInputElement)
	w.sentenceCountElement = w.document.QuerySelector("#sentenceCount").(*dom.HTMLInputElement)
	w.generateSentencesElement = w.document.QuerySelector("#generateSentences").(*dom.HTMLInputElement)
	w.forbidDuplicatesElement = w.document.QuerySelector("#forbidDuplicates").(*dom.HTMLInputElement)
	w.forceWordLimitElement = w.document.QuerySelector("#forceWordLimit").(*dom.HTMLInputElement)
	w.sortOutputElement = w.document.QuerySelector("#sortOutput").(*dom.HTMLInputElement)
	w.markSyllablesElement = w.document.QuerySelector("#markSyllables").(*dom.HTMLInputElement)
	w.applyRejectionsElement = w.document.QuerySelector("#applyRejections").(*dom.HTMLInputElement)
	w.applyReplacementsElement = w.document.QuerySelector("#applyReplacements").(*dom.HTMLInputElement)
	w.copyButtonElement = w.document.QuerySelector("#copyButton").(*dom.HTMLButtonElement)

	w.generatedAlertElement = w.document.QuerySelector("#generatedAlert").(*dom.HTMLDivElement)
	w.duplicateAlertElement = w.document.QuerySelector("#duplicateAlert").(*dom.HTMLDivElement)
	w.rejectedAlertElement = w.document.QuerySelector("#rejectedAlert").(*dom.HTMLDivElement)
	w.replacedAlertElement = w.document.QuerySelector("#replacedAlert").(*dom.HTMLDivElement)

	w.examplePageElements = nil
	if sampleDropdownElement := w.document.QuerySelector("#samples"); sampleDropdownElement != nil {
		w.examplePageElements = &examplePageElements{
			sampleDropdownElement: sampleDropdownElement.(*dom.HTMLSelectElement),
		}
	}
}

func (w *Web) setEventListeners() {
	w.submitButtonElement.AddEventListener("click", false, w.onSubmit)
	w.generateSentencesElement.AddEventListener("click", false, func(event dom.Event) {
		if w.generateSentencesElement.Checked() {
			w.forbidDuplicatesElement.SetDisabled(true)
			w.forceWordLimitElement.SetDisabled(true)
			w.markSyllablesElement.SetDisabled(true)
			w.sortOutputElement.SetDisabled(true)

			w.wordCountElement.SetDisabled(true)
			w.sentenceCountElement.SetDisabled(false)
		} else {
			w.forbidDuplicatesElement.SetDisabled(false)
			w.forceWordLimitElement.SetDisabled(false)
			w.markSyllablesElement.SetDisabled(false)
			w.sortOutputElement.SetDisabled(false)

			w.wordCountElement.SetDisabled(false)
			w.sentenceCountElement.SetDisabled(true)
		}
	})
	// set up copy-output button
	w.copyButtonElement.AddEventListener("click", false, func(event dom.Event) {
		js.Global().Get("window").Get("navigator").Get("clipboard").Call("writeText", w.outputTextElement.Value())
	})
}

func (w *Web) onSubmit(event dom.Event) {
	// get the values of the various options
	w.getOptions()

	// if this is the example page, load the selected code
	if w.examplePageElements != nil {
		selectedExampleValue := w.sampleDropdownElement.Value()
		if selectedExampleValue == "nothing" {
			util.LogError("no example selected", selectedExampleValue)
			w.Evaluator.AddErrors(fmt.Errorf("no example selected"))
			return
		}
		selectedExample, ok := sample.ExampleToFilename[selectedExampleValue]
		if !ok {
			util.LogError("invalid example selection", selectedExampleValue)
			return
		}

		file, err := sample.Examples.Open(selectedExample)
		if err != nil {
			util.LogError("failed to open example", err)
			w.Evaluator.AddErrors(fmt.Errorf("failed to open example: %w", err))
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			util.LogError("failed to read example", err)
			w.Evaluator.AddErrors(fmt.Errorf("failed to read example: %w", err))
			return
		}

		w.inputTextElement.SetValue(string(data))
	}

	w.Evaluator = eval.New(w.getOptions())
	words, sep, sentences := w.Evaluator.Run()
	defer func() {
		if len(w.Evaluator.Errors) > 0 {
			w.displayErrors(w.Evaluator.Errors)
		}
		w.Evaluator.ClearErrors()
	}()

	if w.Evaluator.Options.GenerateSentences {
		w.displaySentences(sentences)
	} else {
		w.displayWords(words, sep)
	}
}

func (w *Web) getOptions() (opts eval.Options) {
	opts.MinSylCount = int(w.minSylCountElement.ValueAsNumber())
	opts.MaxSylCount = int(w.maxSylCountElement.ValueAsNumber())
	opts.WordCount = int(w.wordCountElement.ValueAsNumber())
	opts.SentenceCount = int(w.sentenceCountElement.ValueAsNumber())

	// handle minSylCount being larger than maxSylCount
	if opts.MinSylCount > opts.MaxSylCount {
		opts.MaxSylCount = opts.MinSylCount
		w.maxSylCountElement.SetValue(strconv.Itoa(opts.MinSylCount))
	}

	opts.ForbidDuplicates = w.forbidDuplicatesElement.Checked()
	opts.ForceWordLimit = w.forceWordLimitElement.Checked()
	opts.SortOutput = w.sortOutputElement.Checked()
	opts.MarkSyllables = w.markSyllablesElement.Checked()
	opts.ApplyRejections = w.applyRejectionsElement.Checked()
	opts.ApplyReplacements = w.applyReplacementsElement.Checked()
	opts.GenerateSentences = w.generateSentencesElement.Checked()

	return opts
}

func (w *Web) displayWords(words []eval.Word, syllableSep string) {
	wordStrings := []string{}
	text := ""
	for _, word := range words {
		wordStrings = append(wordStrings, strings.Join(word.Syllables, syllableSep))
	}
	text += strings.Join(wordStrings, "\n")
	w.outputTextElement.SetValue(text)
	w.updateAlerts()
}

func (w *Web) displaySentences(sentences []string) {
	text := strings.Join(sentences, " ")
	w.outputTextElement.SetValue(text)
	w.updateAlerts()
}

func (w *Web) displayErrors(es []error) {
	errs := errors.Join(es...)
	util.LogError(errs)
	w.outputTextElement.SetValue(errs.Error())
}

func (w *Web) updateAlerts() {
	w.generatedAlertElement.SetInnerHTML(fmt.Sprintf("generated %d words", w.Evaluator.GeneratedCount))
	w.duplicateAlertElement.SetInnerHTML(fmt.Sprintf("removed %d duplicates", w.Evaluator.DuplicateCount))
	w.rejectedAlertElement.SetInnerHTML(fmt.Sprintf("rejected %d words", w.Evaluator.RejectedCount))
	w.replacedAlertElement.SetInnerHTML(fmt.Sprintf("replaced %d words", w.Evaluator.ReplacedCount))
}
