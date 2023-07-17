import { CategoryListing } from "../category/Category.js";
import { Syllable } from "../syllable/parser.js";
export default class Reject {
    matchWordStart: boolean;
    matchWordEnd: boolean;
    matchSylStart: boolean;
    matchSylEnd: boolean;
    rejectSyllable: Syllable;
    constructor(rejection: string, categories: CategoryListing);
    toRegex(): RegExp;
}
export { Reject };
