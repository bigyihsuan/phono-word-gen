export declare class SyllableExpr {
}
export declare class Token {
    lexeme: string;
    startingIndex: number;
    endingIndex: number;
    constructor(lexeme: string, startingIndex: number, endingIndex: number);
    toString(): string;
}
export declare function tokenizeSyllable(line: string): Token[];
