package eval

import "slices"

func (e *Evaluator) rejectWords(words []Word) []Word {
	if !e.ApplyRejections {
		return words
	}

	keptWords := []Word{}

	for i, word := range words {
		w, _ := word.Join()

		matchesWordLevel := len(e.wordRejections.String()) > 0 && e.wordRejections.MatchString(w)

		matchesSyllableLevel := false
		if len(e.syllableRejections.String()) > 0 {
			if slices.ContainsFunc(word.Syllables, e.syllableRejections.MatchString) {
				matchesSyllableLevel = true
			}
		}

		matchesGeneral := len(e.generalRejections.String()) > 0 && e.generalRejections.MatchString(w)

		if !matchesWordLevel && !matchesSyllableLevel && !matchesGeneral {
			keptWords = append(keptWords, words[i])
		} else {
			e.RejectedCount++
		}
	}
	return keptWords
}
