# syntax

## Phonemes

separated by spaces

```txt
a b c あ バ ば ɑ θ ɪ any length phonemes
```

## Categories

```py
name = phonemes
```

category definitions end with a newline

category names can only be regex `[a-zA-Z0-9_]+` (alphanum + underscore)

## Syllable Structure

categories are used with preceding `$`

anything that's not a special character, space, or a category is a literal character

```php
($consonant)$vowel[$N x $P]
```

### Syllable Component Modifiers

exactly 1, optional, select 1, random weight

random weight is range from 0 (never) and 1 (always)

random weight can only be applied to optional and selection

if sum of random weight in selection < 1, rest of the chance is filled with "nothing"

```php
$V                          # exactly 1
($C)$V                      # opional
($C ($R))$V                 # nexted optional, makes V, CV, CRV
[$V $C]                     # selection options are separated by spaces
[$V ($C ([$R $X]))$V] $agga # can nest all of these within each other
$C($R)*0.33 $V              # random weight: 0.33 chance to make a R
[$nasal*0.20 $stop*0.30]$V  # random weight: 0.20 chance for nasal, 0.30 chance for stop, **0.50 chance of nothing**    
```
