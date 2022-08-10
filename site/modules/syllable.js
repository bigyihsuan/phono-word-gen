export class SyllableExpr {
}
// syllable lexer state
var SLS;
(function (SLS) {
    SLS[SLS["Start"] = 0] = "Start";
    SLS[SLS["InCategory"] = 1] = "InCategory";
    SLS[SLS["InRaw"] = 2] = "InRaw";
    SLS[SLS["InWeight"] = 3] = "InWeight";
})(SLS || (SLS = {}));
export class Token {
    constructor(lexeme, startingIndex, endingIndex) {
        this.lexeme = lexeme;
        this.startingIndex = startingIndex;
        this.endingIndex = endingIndex;
    }
    toString() {
        return `{${this.constructor.name} ${this.lexeme} @ (${this.startingIndex},${this.endingIndex})}`;
    }
}
class RawComponent extends Token {
}
class Category extends Token {
}
class Lparen extends Token {
}
class Rparen extends Token {
}
class Lbracket extends Token {
}
class Rbracket extends Token {
}
class Comma extends Token {
}
class Star extends Token {
}
class Weight extends Token {
}
const NAME_END = "$*,[]()";
export function tokenizeSyllable(line) {
    // tokenize syllable string
    let sylLine = line.replaceAll("syllable:", "").trim(); // delete the `syllable:` directive
    let state = SLS.Start;
    let tokens = [];
    let idx = 0;
    let lexeme = "";
    let startingIndex = 0;
    while (idx < sylLine.length) {
        let char = sylLine[idx];
        switch (state) {
            case SLS.Start: {
                startingIndex = idx;
                switch (char) {
                    case '(': {
                        idx++;
                        tokens.push(new Lparen(char, startingIndex, idx));
                        break;
                    }
                    case ')': {
                        idx++;
                        tokens.push(new Rparen(char, startingIndex, idx));
                        break;
                    }
                    case '[': {
                        idx++;
                        tokens.push(new Lbracket(char, startingIndex, idx));
                        break;
                    }
                    case ']': {
                        idx++;
                        tokens.push(new Rbracket(char, startingIndex, idx));
                        break;
                    }
                    case '*': {
                        idx++;
                        tokens.push(new Star(char, startingIndex, idx));
                        break;
                    }
                    case ',': {
                        idx++;
                        tokens.push(new Comma(char, startingIndex, idx));
                        break;
                    }
                    case '$': {
                        state = SLS.InCategory;
                        lexeme += char;
                        idx++;
                        break;
                    }
                    case '0': {
                        state = SLS.InWeight;
                        lexeme += char;
                        idx++;
                        break;
                    }
                    default: {
                        state = SLS.InRaw;
                        lexeme += char;
                        idx++;
                        break;
                    }
                }
                break;
            }
            case SLS.InRaw: {
                if (NAME_END.includes(char)) {
                    tokens.push(new RawComponent(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = "";
                }
                else {
                    lexeme += char;
                    idx++;
                }
                break;
            }
            case SLS.InCategory: {
                if (NAME_END.includes(char)) {
                    tokens.push(new Category(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = "";
                }
                else {
                    lexeme += char;
                    idx++;
                }
                break;
            }
            case SLS.InWeight: {
                if (NAME_END.includes(char)) {
                    tokens.push(new Weight(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = "";
                }
                else {
                    lexeme += char;
                    idx++;
                }
                break;
            }
        }
    }
    return tokens;
}
//# sourceMappingURL=syllable.js.map