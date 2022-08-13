import { EvaluableComponent } from "./EvaluableComponent.js";
import { RandomlyChoosable } from "./RandomlyChoosable.js";
export default class OptionalComponent implements EvaluableComponent, RandomlyChoosable {
    component: EvaluableComponent;
    weight: number;
    constructor(component: EvaluableComponent, weight?: number);
    getRandomChoice(): string;
    evaluate(): string;
    toString(): string;
}
export { OptionalComponent };
