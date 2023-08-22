# Specification

- [Specification](#specification)
  - [Lexing](#lexing)
  - [Parsing](#parsing)
    - [Common Tokens](#common-tokens)
    - [Components](#components)
    - [Category Definitions](#category-definitions)
    - [Syllables](#syllables)
    - [Rejections](#rejections)
    - [Replacements](#replacements)
    - [Letters](#letters)
  - [Evaluating](#evaluating)
    - [Preparation](#preparation)
    - [Word Generation](#word-generation)
      - [Weights](#weights)
      - [Categories](#categories)
      - [Syllables](#syllables-1)
      - [Words](#words)
  - [Presentation](#presentation)

## Lexing

The lexer should allow raw source in 2 forms:

- All directives (categories, syllables, replacements, rejections, letters) in one file/text block
- Each of the aforementioned in separate files/text blocks

## Parsing

5 directives, depending on starting token:

| Token                 | Description         | Multiple Allowed? |
| --------------------- | ------------------- | :---------------: |
| `[any category name]` | Category definition |         ✅         |
| `syllable`            | Syllable definition |         ✅         |
| `reject`              | Rejection rule      |         ✅         |
| `replace`             | Replacement rule    |         ✅         |
| `letters`             | Sorting order       |         ❌         |

### Common Tokens

The following grammar is in common for all rules below:

```ebnf
line-ending    = "\n" | ";" ;
with-weight    = "*" weight ;
weight         = [1-9]+[0-9]* ; # any positive decimal integer
comment        = "#" .* "\n" ;
context-prefix = "^" | "@" | "!";
context-suffix = "\" | "&" ;
```

### Components

Components are the base unit of all directives.
They are used in all directives (except for `letters`).
All components can be weighted.

```ebnf
phoneme       = "any non-space text" ;
reference     = "$" category-name ;
category-name = "any non-space text" ;
weighted-component  = syllable-component with-weight? ;
weighted-components = syllable-components with-weight? ;
```

### Category Definitions

Defines a category. One category per line, or multiple when separated by semicolons.

```ebnf
category-definition = category-name "=" space-separated-elements line-ending ;
space-separated-elements = weighted-category-element (" " weighted-category-element)* ;
weighted-category-element = category-element with-weight? ;
category-element = phoneme | reference ;
```

### Syllables

Defines all possible syllables in the language. Multiple syllable directives per file is allowed.

```ebnf
syllable-definition = "syllable" ":" syllable-components line-ending ;
syllable-components = syllable-component+ ;
syllable-component  = phoneme | reference | grouping | selection | optional ;
grouping  = "{" "}" | "{" syllable-components (" "? syllable-components)* "}" ;
optional  = "(" syllable-components ")" with-weight? ;
selection = "[" "]" | "[" selection-elements "]" ;

selection-elements = selection-element ("," selection-element)*
selection-element  = syllable-components with-weight?;
```

### Rejections

Defines a rejection rule. Multiple per file is allowed.
They are of the form `reject: ...` where `...` is any syllable component.

```ebnf
rejection-definition = "reject" ":" rejection-elements line-ending ;
rejection-elements   = rejection-element ("|" rejection-element)? ;
rejection-element    = context-prefix? syllable-components context-suffix? ;
```

### Replacements

```ebnf
replacement-definition = "replace" ":" source? ">" replacement? "/" replace-condition ("//" replace-exception)? line-ending ;
source      = (reference | phoneme)+ ;
replacement = phoneme+ ;
condition   = env ;
exception   = env ;
env         = (context-prefix? syllable-components)? "_" (syllable-components context-suffix?)?
```

### Letters

```ebnf
letters-definition = "letters" ":" letters line-ending ;
letters = phoneme (" " phoneme)* ;
```

## Evaluating

### Preparation

1. Create categories
2. Create syllables
3. Create rejection rules
4. Create replacement rules
5. Create letter sorting order

### Word Generation

1. Pick the syllable count
2. Generate that many syllables
3. Apply rejection rules
4. Apply replacement rules
5. Generate the word's letterization

#### Weights

Weights on components represent how much more probable it is to select it compared to other components.
Components that have the same weight have identical probabilities to be selected.
Based on the semantics of the [weightedrandom](https://pkg.go.dev/github.com/mroth/weightedrand/v2) package.

#### Categories

Categories should assign weights to every phoneme during preparation.
By default, every phoneme has a weight of 1.
Weighted phonemes have their written weight.

When `Get` is called, it should return a phoneme from its colleciton of phonemes.
If the element that was gotten is a reference, it should get from the reference.
Recursive/looped references (`C = $A; A = $C`) are not allowed.
This is checked during preperation, after the categories are generated.

#### Syllables

Syllables should be prepared once, and run many times (assuming the source input doesn't change).
When a syllable's `Get` is called, call `Get` for each component, and turn itself into a string.

#### Words

Words are groupings of multiple syllables.
The number of syllables in a word is dependent on the min/max syllable count options.
The number of words generated is determined by the word count, "allow duplicates", and "force generate to word limit" options.

## Presentation

1. If on, remove duplicates
2. If on, generate again until the word count is reached
3. If on, sort the words based on letters
4. If on, display the word split into syllables
5. Otherwise, display the word as one word