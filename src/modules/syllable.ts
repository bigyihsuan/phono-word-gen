export class SyllableExpr { }

export class Token {
    lexeme: string;
    startingIndex: number;
    endingIndex: number;
    tokenType: string = "";

    constructor(lexeme: string, startingIndex: number, endingIndex: number) {
        this.lexeme = lexeme
        this.startingIndex = startingIndex
        this.endingIndex = endingIndex
    }

    toString(): string {
        return `{${this.constructor.name} ${this.lexeme} @ (${this.startingIndex},${this.endingIndex})}`
    }
}

class RawComponentToken extends Token { tokenType = "RawComponentToken"; }
class CategoryToken extends Token { tokenType = "CategoryToken"; }
class LparenToken extends Token { tokenType = "LparenToken"; }
class RparenToken extends Token { tokenType = "RparenToken"; }
class LbracketToken extends Token { tokenType = "LbracketToken"; }
class RbracketToken extends Token { tokenType = "RbracketToken"; }
class CommaToken extends Token { tokenType = "CommaToken"; }
class StarToken extends Token { tokenType = "StarToken"; }
class WeightToken extends Token { tokenType = "WeightToken"; }


const NAME_END = "$*,[]()";

// syllable lexer state
enum SLS {
    Start,
    InCategory, // `$name`
    InRaw,      // `name`
    InWeight,   // `*0.1234`
}

export function tokenizeSyllable(line: string): Token[] {
    // tokenize syllable string
    let sylLine = line.replaceAll("syllable:", "").trim(); // delete the `syllable:` directive
    let state = SLS.Start;

    let tokens: Token[] = []

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
                        tokens.push(new LparenToken(char, startingIndex, idx));
                        break;
                    }
                    case ')': {
                        idx++;
                        tokens.push(new RparenToken(char, startingIndex, idx));
                        break;
                    }
                    case '[': {
                        idx++;
                        tokens.push(new LbracketToken(char, startingIndex, idx));
                        break;
                    }
                    case ']': {
                        idx++;
                        tokens.push(new RbracketToken(char, startingIndex, idx));
                        break;
                    }
                    case '*': {
                        idx++;
                        tokens.push(new StarToken(char, startingIndex, idx));
                        break;
                    }
                    case ',': {
                        idx++;
                        tokens.push(new CommaToken(char, startingIndex, idx));
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
                        break
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
                    tokens.push(new RawComponentToken(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = ""
                } else {
                    lexeme += char;
                    idx++;
                }
                break;
            }
            case SLS.InCategory: {
                if (NAME_END.includes(char)) {
                    tokens.push(new CategoryToken(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = ""
                } else {
                    lexeme += char;
                    idx++;
                }
                break;
            }
            case SLS.InWeight: {
                if (NAME_END.includes(char)) {
                    tokens.push(new WeightToken(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = ""
                } else {
                    lexeme += char;
                    idx++;
                }
                break;
            }
        }
    }

    return tokens;
}
