import { CategoryListing } from "../category/Category.js";
import { Syllable } from "../syllable/parser.js";
export default class Replacement {
    source: Syllable;
    sourceString: string;
    substitute: Syllable;
    substituteString: string;
    conditionString: string;
    rule: RegExp;
    constructor(replStr: string, categories: CategoryListing);
    matches(word: string): boolean;
    replace(word: string): string;
    apply(word: string): {
        result: string;
        couldApply: boolean;
    };
    toString(): string;
}
export { Replacement };
