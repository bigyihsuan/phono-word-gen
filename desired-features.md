# Desired Features

## Minimum

- `DONE` Arbitrary-length phonemes
- `DONE` Phoneme categories
- `DONE` Optional, Selection (`(a)`, `[a,b,]`)
- `DONE` Randomness Control: Manual weighting of components in selections
- `DONE` Syllable-based word generator
- `DONE` Can specify minimum and maximum number of syllables per word.

## After Minimum

- `DONE` Randomness Control: Manual weighting per-phoneme in categories
- `DONE` Treating categories like sets, and allowing addition of categories
- Word filtration, substitution, modification (Lexurgy-style rewrite rules?)
  - phoneme-based
  - category-based
- Cluster tables
  - phoneme-based only

## Nice to Haves

- Randomness control: Based on worldwide phoneme frequency (via [PHOIBLE data](https://phoible.org/))

## TODOs

- change parser to allow for weights on selection options when not grouped
