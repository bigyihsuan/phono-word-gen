import { EvaluableComponent } from "./EvaluableComponent.js";
import { RandomlyChoosable } from "./RandomlyChoosable.js";

export default class OptionalComponent implements EvaluableComponent, RandomlyChoosable {
    component: EvaluableComponent;

    weight: number;

    // default weight is 50/50
    constructor(component: EvaluableComponent, weight: number = 0.5) {
        this.component = component;
        this.weight = weight;
    }

    getRandomChoice(): string {
        return Math.random() < this.weight ? this.component.evaluate() : "";
    }

    evaluate(): string {
        return this.getRandomChoice();
    }
}

export { OptionalComponent };
