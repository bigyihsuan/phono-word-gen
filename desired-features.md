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
- `DONE` Word filtration (`reject:`)
  - `DONE` Start of word `<^$C>`, end of word `<$C;>`
  - Start of syllable `<@$C>`, end of syllable `<$C&>`
  - Negation? (`<!$C>` = "everything but the things in `C`)
- Word substitution/modification (Lexurgy-style rewrite rules?)
- Cluster tables
  - phoneme-based only

## Nice to Haves

- Randomness control: Based on worldwide phoneme frequency (via [PHOIBLE data](https://phoible.org/))

## TODOs
