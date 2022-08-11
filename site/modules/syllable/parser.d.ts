import { CategoryListing } from "./../category.js";
import { Token } from "./token.js";
interface EvaluableComponent {
    evaluate(): string;
}
declare class ParseError {
    reason: string;
    token: Token;
    sylStr: string;
    constructor(reason: string, token: Token, sylStr: string);
    toString(): string;
}
declare class Syllable implements EvaluableComponent {
    components: SyllableExpr[];
    constructor(components: SyllableExpr[]);
    evaluate(): string;
}
declare class SyllableExpr implements EvaluableComponent {
    component: EvaluableComponent;
    constructor(component: EvaluableComponent);
    evaluate(): string;
}
declare function parseSyllable(tokens: Token[], categories: CategoryListing, sylStr: string): Syllable | ParseError;
export { Syllable, parseSyllable, ParseError, };
