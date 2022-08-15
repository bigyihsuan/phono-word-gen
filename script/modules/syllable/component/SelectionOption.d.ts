import { EvaluableComponent } from "./EvaluableComponent.js";
export default class SelectionOption implements EvaluableComponent {
    component: EvaluableComponent;
    weight: number;
    constructor(component: EvaluableComponent, weight: number);
    toString(): string;
    evaluate(): string;
}
export { SelectionOption };
