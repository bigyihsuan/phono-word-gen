package eval

import (
	"slices"
	"sort"
	"strings"
)

func (e *Evaluator) sort(words []Word) []Word {
	// letter-based sorting
	if len(e.letters) > 0 {
		sort.Slice(words, func(left, right int) bool {
			// letterize words
			// join into a single string
			l := strings.Join(words[left].Syllables, "")
			r := strings.Join(words[right].Syllables, "")
			// find all (known) letters
			leftLetters := e.letterRegexp.FindAllString(l, -1)
			rightLetters := e.letterRegexp.FindAllString(r, -1)
			// for each letter found, find the index of that letter in the letter directive
			leftIndexes := []int{}
			rightIndexes := []int{}
			for _, letter := range leftLetters {
				leftIndexes = append(leftIndexes, slices.Index(e.letters, letter))
			}
			for _, letter := range rightLetters {
				rightIndexes = append(rightIndexes, slices.Index(e.letters, letter))
			}
			minLen := min(len(leftIndexes), len(rightIndexes))

			for i := 0; i < minLen; i++ {
				if leftIndexes[i] < rightIndexes[i] {
					return true
				}
				if leftIndexes[i] > rightIndexes[i] {
					return false
				}
			}
			if len(leftIndexes) < len(rightIndexes) {
				return true
			}
			if len(leftIndexes) > len(rightIndexes) {
				return false
			}
			return false
		})
	} else {
		sort.Slice(words, func(i, j int) bool {
			a, b := words[i], words[j]
			as, bs := strings.Join(a.Syllables, ""), strings.Join(b.Syllables, "")
			less := as < bs
			// if less {
			// 	words[i], words[j] = words[j], words[i]
			// }
			return less
		})
	}
	return words
}
