# my random word generator rant thing posted to the Conlangs Discord Network's `#tools-and-documentation`

(with some minimal formatting edits)

man i have like 4 different random word generators bookmarked and all of them have features missing for my "ideal" random word generator

- [my random word generator rant thing posted to the Conlangs Discord Network's `#tools-and-documentation`](#my-random-word-generator-rant-thing-posted-to-the-conlangs-discord-networks-tools-and-documentation)
  - [gen](#gen)
  - [awkwords](#awkwords)
  - [gengo](#gengo)
  - [lexifier](#lexifier)
  - [MY IDEAL WORD GENERATOR](#my-ideal-word-generator)

## gen

<http://www.zompist.com/gen.html>

liked

- ???

my problems

- phonemes need to be 1 char each, di-/tri-graphs need to be done with rewrite rules
- ordering of phonemes matters (with dropoff != "equiprobable")
- need to list *ALL* variants of syllables
- can't set to no monosyllables

## awkwords

<http://akana.conlang.org/tools/awkwords/>

liked

- arbitrary length phonemes
- optional and "pick one" for syllable patterns
- can set probabilities

my problems

- max of 26 categories
- can't set syllable structure and generate words of `N` syllables, have to manually copy-paste syllable structure multiple times

## gengo

<https://collinbrennan.github.io/GenGo/index.html>

liked

- arbitrary length phonemes
- optional and "pick one" for syllable patterns
- can set probabilities
- rewrites
- specify syllable structure and generate words of `N` to `M`-syllable words

my problems

- incomplete, according to docs (and the github commit history)
- weird comment syntax

## lexifier

<https://lingweenie.org/conlang/lexifer-app.html>

liked

- arbitrary length phonemes
- rejecting, fitering
- syllable macros
- optional syllable components
- cluster fields

my problems

- no "pick one" for words
- works on word structure rather than syllable structure
- need to list many combos of syllables
- macros are only 1 character long
- no classes in classes
- order of phonemes matters

## MY IDEAL WORD GENERATOR

- arbitrary length phonemes (awk, gengo, lex)
- rejecting, filtering, rewrites (gen, awk, gengo, lex)
- optional, "pick one", repeated syllable components (awk, gengo, ½lex)
- specify syllable, generate `N-M`-syllable words (gengo)
- pure random (½gen, awk, gengo), or based on worldwide phoneme frequency (none of these, yet)
