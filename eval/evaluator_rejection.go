package eval

func (e *Evaluator) rejectWords(words []Word) []Word {
	if !e.applyRejections {
		return words
	}

	keptWords := []Word{}

	for i, word := range words {
		w, _ := word.Join()

		matchesWordLevel := len(e.wordRejections.String()) > 0 && e.wordRejections.MatchString(w)

		matchesSyllableLevel := false
		if len(e.syllableRejections.String()) > 0 {
			for _, syl := range word.Syllables {
				if e.syllableRejections.MatchString(syl) {
					matchesSyllableLevel = true
					break
				}
			}
		}

		matchesGeneral := len(e.generalRejections.String()) > 0 && e.generalRejections.MatchString(w)

		if !matchesWordLevel && !matchesSyllableLevel && !matchesGeneral {
			keptWords = append(keptWords, words[i])
		} else {
			e.rejectedCount++
		}
	}
	return keptWords
}
