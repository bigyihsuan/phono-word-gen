export class ParseError {
    constructor(reason, token, sylStr) {
        this.reason = reason;
        this.token = token;
        this.sylStr = sylStr;
    }
    toString() {
        const endIndex = this.token.endingIndex - this.token.startingIndex - 1;
        const locationstr = "^".padStart(this.token.startingIndex + 1, " ") + (endIndex > 0 ? "^".padStart(endIndex, "-") : "");
        return [this.reason, this.sylStr, locationstr].join("\n");
    }
}
//# sourceMappingURL=ParseError.js.map