import { CategoryListing } from "../category.js";
import { ParseError } from "./ParseError.js";
import { Syllable } from "./component/Syllable.js";
import { Token } from "./token.js";
declare function parseSyllable(tokens: Token[], categories: CategoryListing, sylStr: string): Syllable | ParseError;
export { Syllable, parseSyllable, ParseError, };
