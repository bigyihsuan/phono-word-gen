import {
    Token,
    LparenToken,
    RparenToken,
    LbracketToken,
    RbracketToken,
    StarToken,
    CommaToken,
    RawComponentToken,
    CategoryToken,
    WeightToken,
    LcurlyToken,
    RcurlyToken,
} from "./token.js";

const NAME_END = "$*,[](){} \n";

// syllable lexer state
enum SLS {
    Start = "Start",
    InCategory = "InCategory", // `$name`
    InRaw = "InRaw", // `name`
    InWeight = "InWeight", // `*0.1234`
}

export default function tokenizeSyllable(line: string): Token[] {
    // tokenize syllable string
    const sylLine = `${line.trim()}\n`; // delete the `syllable:` directive
    let state = SLS.Start;

    const tokens: Token[] = [];

    let idx = 0;
    let lexeme = "";
    let startingIndex = 0;
    while (idx < sylLine.length) {
        const char = sylLine[idx];
        switch (state) {
            case SLS.Start: {
                startingIndex = idx;
                switch (char) {
                    case "(": {
                        idx += 1;
                        tokens.push(new LparenToken(char, startingIndex, idx));
                        break;
                    }
                    case ")": {
                        idx += 1;
                        tokens.push(new RparenToken(char, startingIndex, idx));
                        break;
                    }
                    case "[": {
                        idx += 1;
                        tokens.push(new LbracketToken(char, startingIndex, idx));
                        break;
                    }
                    case "]": {
                        idx += 1;
                        tokens.push(new RbracketToken(char, startingIndex, idx));
                        break;
                    }
                    case "{": {
                        idx += 1;
                        tokens.push(new LcurlyToken(char, startingIndex, idx));
                        break;
                    }
                    case "}": {
                        idx += 1;
                        tokens.push(new RcurlyToken(char, startingIndex, idx));
                        break;
                    }
                    case "*": {
                        idx += 1;
                        tokens.push(new StarToken(char, startingIndex, idx));
                        break;
                    }
                    case ",": {
                        idx += 1;
                        tokens.push(new CommaToken(char, startingIndex, idx));
                        break;
                    }
                    case "$": {
                        state = SLS.InCategory;
                        lexeme += char;
                        idx += 1;
                        break;
                    }
                    default: {
                        if (char.match(/[0-9]/)) {
                            state = SLS.InWeight;
                            lexeme += char;
                            idx += 1;
                            break;
                        } else if (char.match(/\s/)) {
                            // ignore whitespace
                            idx += 1;
                        } else {
                            state = SLS.InRaw;
                            lexeme += char;
                            idx += 1;
                            break;
                        }
                    }
                }
                break;
            }
            case SLS.InRaw: {
                if (NAME_END.includes(char)) {
                    tokens.push(new RawComponentToken(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = "";
                } else {
                    lexeme += char;
                    idx += 1;
                }
                break;
            }
            case SLS.InCategory: {
                if (NAME_END.includes(char)) {
                    tokens.push(new CategoryToken(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = "";
                } else {
                    lexeme += char;
                    idx += 1;
                }
                break;
            }
            case SLS.InWeight: {
                if (NAME_END.includes(char)) {
                    tokens.push(new WeightToken(lexeme, startingIndex, idx));
                    state = SLS.Start;
                    lexeme = "";
                } else {
                    lexeme += char;
                    idx += 1;
                }
                break;
            }
            default: {
                // nop, state is always SLS
            }
        }
    }
    return tokens;
}

export {
    tokenizeSyllable,
};
