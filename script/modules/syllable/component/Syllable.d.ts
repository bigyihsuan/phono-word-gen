import { EvaluableComponent } from "./EvaluableComponent.js";
import { SyllableExpr } from "./SyllableExpr.js";
export default class Syllable implements EvaluableComponent {
    components: SyllableExpr[];
    constructor(components: SyllableExpr[]);
    evaluate(): string;
    toString(): string;
}
export { Syllable };
