import { IEvaluableComponent } from "./IEvaluableComponent.js";

export default class SyllableExpr implements IEvaluableComponent {
    component: IEvaluableComponent;

    constructor(component: IEvaluableComponent) {
        this.component = component;
    }

    evaluate(): string {
        return this.component.evaluate();
    }

    evaluateAll(): string[] {
        return this.component.evaluateAll();
    }

    toString(): string {
        return this.component.toString();
    }
}

export { SyllableExpr };
