/* eslint-disable max-classes-per-file */
class Token {
    lexeme;
    startingIndex;
    endingIndex;
    constructor(lexeme, startingIndex, endingIndex) {
        this.lexeme = lexeme;
        this.startingIndex = startingIndex;
        this.endingIndex = endingIndex;
    }
    toString() {
        return `{${this.constructor.name} ${this.lexeme} @ (${this.startingIndex},${this.endingIndex})}`;
    }
}
class RawComponentToken extends Token {
}
class CategoryToken extends Token {
}
class LparenToken extends Token {
}
class RparenToken extends Token {
}
class LbracketToken extends Token {
}
class RbracketToken extends Token {
}
class LcurlyToken extends Token {
}
class RcurlyToken extends Token {
}
class CommaToken extends Token {
}
class StarToken extends Token {
}
class WeightToken extends Token {
}
export { Token, RawComponentToken, CategoryToken, LparenToken, RparenToken, LbracketToken, RbracketToken, CommaToken, StarToken, WeightToken, LcurlyToken, RcurlyToken, };
//# sourceMappingURL=token.js.map