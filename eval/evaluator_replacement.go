package eval

import "phono-word-gen/util"

func (e *Evaluator) replaceWords(words []Word) []Word {
	if !e.applyReplacements {
		return words
	}

	replacedWords := []Word{}

	// for each word...
	for _, word := range words {
		w, _ := word.Join()
		// w, syllableIndexes := word.Join()

		// for each possible replacement...
		for _, r := range e.replacements {
			matchesException := false

			// check for exception; if true, don't replace
			exception := r.ExceptionRegexp(e.categories, e.components)
			if exception != nil {
				if r.Exception.IsSyllableLevel() {
					matchesException = true
					for _, syl := range word.Syllables {
						if !exception.MatchString(syl) {
							matchesException = false
							break
						}
					}
				} else {
					matchesException = exception.MatchString(w)
				}
			}
			if matchesException {
				// return early
				replacedWords = append(replacedWords, word)
				continue
			}

			// need to find the match indexes to replace across syllable boundaries
			// match against the word, and get the indexes
			// map word-indexes into syllable-letter-indexes
			//     word[i] => syllables[j][k]
			// this is so that the replacement can span across syllables
			// though, there's too many parens in the generated regexp
			// so need to figure out how to get the start and end index
			// w/out using subgroups

			condition := r.ConditionRegexp(e.categories, e.components)
			matchIndexes := [][]int{}
			if r.Condition.IsSyllableLevel() {
				for _, syl := range word.Syllables {
					match := condition.FindAllStringIndex(syl, -1)
					if match == nil {
						continue
					}
					matchIndexes = append(matchIndexes, match...)
				}
			} else {
				matchIndexes = condition.FindAllStringIndex(w, -1)
			}
			m, _ := util.ToMap(matchIndexes)
			util.Log(word, m)
			util.Log(r.Source.Regexp(e.categories, e.components), r.Replacement)
			util.Log("ConditionRegexp", r.ConditionRegexp(e.categories, e.components))
			if matchIndexes == nil || len(matchIndexes) == 0 {
				// return early
				replacedWords = append(replacedWords, word)
				continue
			}
			// TODO: replacement
			// we have the indexes of the matches in either the joined word or syllables

		}
	}
	return replacedWords
}
