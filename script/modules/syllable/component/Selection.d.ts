import { IEvaluableComponent } from "./IEvaluableComponent.js";
import { IRandomlyChoosable } from "./IRandomlyChoosable.js";
import { SelectionOption } from "./SelectionOption.js";
export default class Selection implements IRandomlyChoosable, IEvaluableComponent {
    options: SelectionOption[];
    weights: number[];
    constructor(options: SelectionOption[]);
    generateWeights(): void;
    getRandomChoice(): string;
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
}
export { Selection };
