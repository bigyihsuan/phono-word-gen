import { EvaluableComponent } from "./EvaluableComponent.js";

export default class SyllableExpr implements EvaluableComponent {
    component: EvaluableComponent;

    constructor(component: EvaluableComponent) {
        this.component = component;
    }

    evaluate(): string {
        return this.component.evaluate();
    }
}

export { SyllableExpr };
