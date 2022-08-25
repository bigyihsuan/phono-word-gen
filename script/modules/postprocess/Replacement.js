import tokenizeSyllable from "../syllable/lexer.js";
import { ParseError, parseSyllable } from "../syllable/parser.js";
import Reject from "./Reject.js";
/** capture groups:
 *
 * 0: the original string,
 * 1: the source syllable,
 * 2: the replacement,
 * 3: the conditions
 */
const replaceRegex = /{(.*?)}\s*>\s*{(.*?)}\s*\/\s*(.*)/g;
const replaceWithExceptionRegex = /{(.*?)}\s*>\s*{(.*?)}\s*\/\s*(.*?)\s*\/\/\s*(.*)/g;
export default class Replacement {
    source;
    sourceString;
    substitute;
    substituteString;
    conditionString;
    exceptionString;
    rule;
    hasException;
    constructor(replStr, categories) {
        let [source, substitute, conditions, exceptions] = ["", "", "", ""];
        let result = [];
        this.hasException = replStr.includes("//");
        if (this.hasException) {
            result = Array.from(replStr.matchAll(replaceWithExceptionRegex));
        }
        else {
            result = Array.from(replStr.matchAll(replaceRegex));
        }
        // console.log({ result, resultArr: [...result] });
        source = String(result[0][1]);
        substitute = String(result[0][2]);
        conditions = String(result[0][3]);
        if (this.hasException) {
            exceptions = String(result[0][4]);
        }
        // source
        const so = parseSyllable(tokenizeSyllable(source), categories, source);
        if (so instanceof ParseError) {
            throw so;
        }
        this.source = so;
        this.sourceString = source;
        // substitution
        const su = parseSyllable(tokenizeSyllable(substitute), categories, substitute);
        if (su instanceof ParseError) {
            throw su;
        }
        this.substitute = su;
        this.substituteString = substitute;
        // condition
        let [left, right] = conditions.split("_", 2);
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
        // exception
        let leftRegExc = "";
        let rightRegExc = "";
        if (this.hasException) {
            [left, right] = exceptions.split("_", 2);
            if (left.length > 0) {
                const leftExc = new Reject(left, categories);
                let s = leftExc.toRegex().source;
                s = s.replace("(?:)", "");
                leftRegExc = s;
            }
            if (right.length > 0) {
                const rightExc = new Reject(right, categories);
                let s = rightExc.toRegex().source;
                s = s.replace("(?:)", "");
                rightRegExc = s;
            }
        }
        this.exceptionString = exceptions;
        // generate rule
        let ruleReg = this.source.toRegex().source;
        if (leftReg.length > 0) {
            ruleReg = `(?<=${leftReg})${ruleReg}`;
        }
        if (leftRegExc.length > 0) {
            ruleReg = `(?<!${leftRegExc})${ruleReg}`;
        }
        if (rightReg.length > 0) {
            ruleReg = `${ruleReg}(?=${rightReg})`;
        }
        if (rightRegExc.length > 0) {
            ruleReg = `${ruleReg}(?!${rightRegExc})`;
        }
        this.rule = new RegExp(ruleReg);
        // console.log({
        //     source, substitute, conditions, exceptions, rule: this.rule,
        // });
    }
    matches(word) {
        return this.rule.test(word);
    }
    replace(word) {
        // TODO? replace with categories after switching to array-based vs. set-based
        return word.replace(this.rule, this.substituteString);
    }
    apply(word) {
        if (this.matches(word)) {
            return { result: this.replace(word), couldApply: true };
        }
        return { result: word, couldApply: false };
    }
    toString() {
        return `${this.sourceString} > ${this.substituteString} / ${this.conditionString} // ${this.exceptionString}`;
    }
}
export { Replacement };
//# sourceMappingURL=Replacement.js.map