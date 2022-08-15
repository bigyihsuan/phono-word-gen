import { IEvaluableComponent } from "./IEvaluableComponent.js";
import { IRandomlyChoosable } from "./IRandomlyChoosable.js";
export default class OptionalComponent implements IEvaluableComponent, IRandomlyChoosable {
    component: IEvaluableComponent;
    weight: number;
    constructor(component: IEvaluableComponent, weight?: number);
    getRandomChoice(): string;
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
}
export { OptionalComponent };
