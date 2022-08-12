import { Token } from "./token.js";
export declare class ParseError {
    reason: string;
    token: Token;
    sylStr: string;
    constructor(reason: string, token: Token, sylStr: string);
    toString(): string;
}
