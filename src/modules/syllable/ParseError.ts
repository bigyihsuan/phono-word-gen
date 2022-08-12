import { Token } from "./token.js";

export class ParseError {
    reason: string;

    token: Token;

    sylStr: string;

    constructor(reason: string, token: Token, sylStr: string) {
        this.reason = reason;
        this.token = token;
        this.sylStr = sylStr;
    }

    toString(): string {
        const endIndex = this.token.endingIndex - this.token.startingIndex - 1;
        const locationstr = "^".padStart(this.token.startingIndex + 1, " ") + (endIndex > 0 ? "^".padStart(endIndex, "-") : "");
        return [this.reason, this.sylStr, locationstr].join("\n");
    }
}
