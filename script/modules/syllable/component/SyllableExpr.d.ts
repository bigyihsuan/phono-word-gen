import { EvaluableComponent } from "./EvaluableComponent.js";
export default class SyllableExpr implements EvaluableComponent {
    component: EvaluableComponent;
    constructor(component: EvaluableComponent);
    evaluate(): string;
    toString(): string;
}
export { SyllableExpr };
