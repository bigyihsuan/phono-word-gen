import { IEvaluableComponent } from "./IEvaluableComponent.js";
import { IRandomlyChoosable } from "./IRandomlyChoosable.js";

export default class OptionalComponent implements IEvaluableComponent, IRandomlyChoosable {
    component: IEvaluableComponent;

    weight: number;

    // default weight is 50/50
    constructor(component: IEvaluableComponent, weight: number = 0.5) {
        this.component = component;
        this.weight = weight;
    }

    getRandomChoice(): string {
        return Math.random() < this.weight ? this.component.evaluate() : "";
    }

    evaluate(): string {
        return this.getRandomChoice();
    }

    evaluateAll(): string[] {
        return ["", ...this.component.evaluateAll()];
    }

    toString(): string {
        return `(${this.component.toString()}:${this.weight})`;
    }
}

export { OptionalComponent };
