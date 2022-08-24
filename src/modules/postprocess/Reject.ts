import { CategoryListing } from "../category/Category.js";
import tokenizeSyllable from "../syllable/lexer.js";
import { ParseError } from "../syllable/ParseError.js";
import { Syllable, parseSyllable } from "../syllable/parser.js";

export default class Reject {
    matchWordStart: boolean;

    matchWordEnd: boolean;

    matchSylStart: boolean;

    matchSylEnd: boolean;

    rejectSyllable: Syllable;

    constructor(rejection: string, categories: CategoryListing) {
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

    toRegex(): RegExp {
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
