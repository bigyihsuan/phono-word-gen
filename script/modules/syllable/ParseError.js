export default class ParseError {
    reason;
    token;
    sylStr;
    withins = [];
    constructor(reason, token, sylStr, sourceFun) {
        this.reason = reason;
        this.token = token;
        this.sylStr = sylStr;
        this.withins.push(sourceFun);
    }
    toString() {
        const endIndex = this.token.endingIndex - this.token.startingIndex - 1;
        const locationstr = "^".padStart(this.token.startingIndex + 1, " ") + (endIndex > 0 ? "^".padStart(endIndex, "-") : "");
        return [this.withins.join("\n     within "), this.reason, this.sylStr, locationstr].join("\n");
    }
    within(funName) {
        this.withins.push(funName);
        return this;
    }
}
export { ParseError };
//# sourceMappingURL=ParseError.js.map