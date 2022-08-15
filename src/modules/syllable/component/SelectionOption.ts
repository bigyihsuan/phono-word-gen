import { IEvaluableComponent } from "./IEvaluableComponent.js";

export default class SelectionOption implements IEvaluableComponent {
    component: IEvaluableComponent;

    weight: number;

    constructor(component: IEvaluableComponent, weight: number) {
        this.component = component;
        this.weight = weight;
    }

    toString(): string {
        return `${this.component.toString()}:${this.weight}`;
    }

    evaluate(): string {
        return this.component.evaluate();
    }

    evaluateAll(): string[] {
        return this.component.evaluateAll();
    }
}

export { SelectionOption };
