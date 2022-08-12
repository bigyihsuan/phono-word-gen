declare class Token {
    lexeme: string;
    startingIndex: number;
    endingIndex: number;
    constructor(lexeme: string, startingIndex: number, endingIndex: number);
    toString(): string;
}
declare class RawComponentToken extends Token {
}
declare class CategoryToken extends Token {
}
declare class LparenToken extends Token {
}
declare class RparenToken extends Token {
}
declare class LbracketToken extends Token {
}
declare class RbracketToken extends Token {
}
declare class LcurlyToken extends Token {
}
declare class RcurlyToken extends Token {
}
declare class CommaToken extends Token {
}
declare class StarToken extends Token {
}
declare class WeightToken extends Token {
}
export { Token, RawComponentToken, CategoryToken, LparenToken, RparenToken, LbracketToken, RbracketToken, CommaToken, StarToken, WeightToken, LcurlyToken, RcurlyToken, };
