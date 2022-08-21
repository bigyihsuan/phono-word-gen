import { Token } from "./token.js";
export default class ParseError {
    reason: string;
    token: Token;
    sylStr: string;
    withins: string[];
    constructor(reason: string, token: Token, sylStr: string, sourceFun: string);
    toString(): string;
    within(funName: string): ParseError;
    appendMessage(message: string): ParseError;
}
export { ParseError };
