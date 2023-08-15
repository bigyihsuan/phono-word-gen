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
  - [Presentation](#presentation)

## Lexing

The lexer should allow raw source in 2 forms:

- All directives (categories, syllables, replacements, rejections, letters) in one file/text block
- Each of the aforementioned in separate files/text blocks

## Parsing

5 directives, depending on starting token:

| Token                 | Description         | Multiple Allowed? |
| --------------------- | ------------------- | :---------------: |
| `syllable`            | Syllable definition |         ✅         |
| `reject`              | Rejection rule      |         ✅         |
| `replace`             | Replacement rule    |         ✅         |
| `letters`             | Sorting order       |         ❌         |
| `[any category name]` | Category definition |         ✅         |

### Common Tokens

The following grammar is in common for all rules below:

```ebnf
line-ending    = "\n" | ";" ;
with-weight    = "*" weight ;
weight         = [0-9]+.[0-9]+ ; # any decimal number
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

Defines all possible syllables in the language.
Multiple syllable directives per file is allowed; all kinds of syllable will be randomly selected equally.
If you do not want uniform randomness for syllables, explicitly wrap them in a weighted selection.

```ebnf
syllable-definition = "syllable" ":" syllable-components line-ending ;
syllable-components = syllable-component (" "? syllable-component)+ ;
syllable-component  = phoneme | reference | grouping | selection | optional ;
grouping  = "{" grouping-elements "}" ;
selection = "[" selection-elements "]" ;
optional  = "(" grouping-elements ")" with-weight? ;

grouping-elements  = syllable-components ;
selection-elements = weighted-components ("," weighted-components)* ;
```

### Rejections

Defines a rejection rule. Multiple per file is allowed.
They are of the form `reject: ...` where `...` is any syllable component.

```ebnf
rejection-definition = "reject" ":" rejection-elements line-ending ;
rejection-elements   = rejection-element ("|" rejection-element)? ;
rejection-element    = context-prefix? category-element+ context-suffix? ;
```

### Replacements

```ebnf
replacement-definition = "replace" ":" source? ">" replacement? "/" replace-condition ("//" replace-exception)? line-ending ;
source      = (reference | phoneme)+ ;
replacement = phoneme ;
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
3. Create replacement rules
4. Create rejection rules
5. Create letter sorting order

### Word Generation

1. Pick the syllable count
2. Generate that many syllables
3. Apply replacement rules
4. Apply rejection rules
5. Generate the word's letterization

## Presentation

1. If on, sort the words based on letters
2. If on, display the word split into syllables
3. Otherwise, display the word as one word