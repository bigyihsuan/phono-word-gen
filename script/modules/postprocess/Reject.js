import tokenizeSyllable from "../syllable/lexer.js";
import { ParseError } from "../syllable/ParseError.js";
import { parseSyllable } from "../syllable/parser.js";
export default class Reject {
    matchWordStart;
    matchWordEnd;
    matchSylStart;
    matchSylEnd;
    rejectSyllable;
    constructor(rejection, categories) {
        this.matchWordStart = rejection.startsWith("^");
        this.matchWordEnd = rejection.endsWith(";");
        this.matchSylStart = rejection.startsWith("@");
        this.matchSylEnd = rejection.endsWith("&");
        const s = parseSyllable(tokenizeSyllable(rejection.replaceAll(/[&^;@]/g, "")), categories, rejection.replaceAll(/[&^;@]/g, ""));
        if (s instanceof ParseError) {
            throw s;
        }
        this.rejectSyllable = s;
    }
    toRegex() {
        let reg = this.rejectSyllable.toRegex().source;
        if (this.matchWordStart) {
            reg = `^${reg}`;
        }
        if (this.matchWordEnd) {
            reg = `${reg}$`;
        }
        return new RegExp(`(${reg})`);
    }
}
export { Reject };
//# sourceMappingURL=Reject.js.map