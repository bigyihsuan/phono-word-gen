# Desired Features

## Minimum

- Arbitrary-length phonemes
- Phoneme categories
- Both phoneme-based and category-based output filtration, substitution, modification
- Optional, "Pick One", repeated syllable components (regex `a*`, `[ab]`, and `a{N,M}`)
- Syllable-based word generator
- Can specify minimum and maximum number of syllables per word.

## Nice to Haves

- Randomness control:
  - Manual specification (per-phoneme, per-syllable component)
  - Based on worldwide phoneme frequency (via [PHOIBLE data](https://phoible.org/))
