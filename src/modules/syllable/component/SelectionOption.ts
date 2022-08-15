import { EvaluableComponent } from "./EvaluableComponent.js";

export default class SelectionOption implements EvaluableComponent {
    component: EvaluableComponent;

    weight: number;

    constructor(component: EvaluableComponent, weight: number) {
        this.component = component;
        this.weight = weight;
    }

    toString(): string {
        return `${this.component.toString()}:${this.weight}`;
    }

    evaluate(): string {
        return this.component.evaluate();
    }
}

export { SelectionOption };
