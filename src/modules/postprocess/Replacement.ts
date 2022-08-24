import { CategoryListing } from "../category/Category.js";
import tokenizeSyllable from "../syllable/lexer.js";
import { ParseError, parseSyllable, Syllable } from "../syllable/parser.js";
import Reject from "./Reject.js";

/** capture groups:
 *
 * 0: the original string,
 * 1: the source syllable,
 * 2: the replacement,
 * 3: the conditions
 */
const replaceRegex = /{(.+)}\s*>\s*{(.+)}\s*\/\s*(.+)/g;

export default class Replacement {
    source: Syllable;

    sourceString: string;

    substitute: Syllable;

    substituteString: string;

    conditionString: string;

    rule: RegExp;

    // sub: string;

    constructor(replStr: string, categories: CategoryListing) {
        const result = replStr.matchAll(replaceRegex);

        const [, source, substitute, conditions] = result.next().value;
        // console.log({
        //     source, substitute, conditions,
        // });

        const so = parseSyllable(tokenizeSyllable(source), categories, source);
        if (so instanceof ParseError) {
            throw so;
        }
        this.source = so;
        this.sourceString = source;

        const su = parseSyllable(tokenizeSyllable(substitute), categories, substitute);
        if (su instanceof ParseError) {
            throw su;
        }
        this.substitute = su;
        this.substituteString = substitute;

        const [left, right] = conditions.split("_", 2);
        let leftReg = "";
        let rightReg = "";

        if (left.length > 0) {
            const leftSyl = new Reject(left, categories);
            let s = leftSyl.toRegex().source;
            s = s.replace("(?:)", "");
            leftReg = s;
        }
        if (right.length > 0) {
            const rightSyl = new Reject(right, categories);
            let s = rightSyl.toRegex().source;
            s = s.replace("(?:)", "");
            rightReg = s;
        }
        this.conditionString = conditions;

        this.rule = new RegExp(`(?<=${leftReg})${this.source.toRegex().source}(?=${rightReg})`);
    }

    matches(word: string): boolean {
        return this.rule.test(word);
    }

    replace(word: string): string {
        return word.replace(this.rule, this.substituteString);
    }

    apply(word: string): { result: string, couldApply: boolean } {
        if (this.matches(word)) {
            return { result: this.replace(word), couldApply: true };
        }
        return { result: word, couldApply: false };
    }

    toString(): string {
        return `${this.sourceString} > ${this.substituteString} / ${this.conditionString}`;
    }
}

export { Replacement };
