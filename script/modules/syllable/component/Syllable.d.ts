import { EvaluableComponent } from "./EvaluableComponent.js";
export default class Syllable implements EvaluableComponent {
    components: EvaluableComponent[];
    constructor(components: EvaluableComponent[]);
    evaluate(): string;
    toString(): string;
}
export { Syllable };
