import { EvaluableComponent } from "./EvaluableComponent.js";
import { RandomlyChoosable } from "./RandomlyChoosable.js";
import { SelectionOption } from "./SelectionOption.js";
export default class Selection implements RandomlyChoosable, EvaluableComponent {
    options: SelectionOption[];
    weights: number[];
    constructor(options: SelectionOption[]);
    generateWeights(): void;
    getRandomChoice(): string;
    evaluate(): string;
}
export { Selection };
