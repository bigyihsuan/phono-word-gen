import { Token } from "./token.js";

export default class ParseError {
    reason: string;

    token: Token;

    sylStr: string;

    withins: string[] = [];

    constructor(reason: string, token: Token, sylStr: string, sourceFun: string) {
        this.reason = reason;
        this.token = token;
        this.sylStr = sylStr;
        this.withins.push(sourceFun);
    }

    toString(): string {
        const endIndex = this.token.endingIndex - this.token.startingIndex - 1;
        const locationstr = "^".padStart(this.token.startingIndex + 1, " ") + (endIndex > 0 ? "^".padStart(endIndex, "-") : "");
        return [this.withins.join("\n     within "), this.reason, this.sylStr, locationstr].join("\n");
    }

    within(funName: string): ParseError {
        this.withins.push(funName);
        return this;
    }

    appendMessage(message: string): ParseError {
        this.reason += message;
        return this;
    }
}

export { ParseError };
