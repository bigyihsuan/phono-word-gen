import { IEvaluableComponent } from "./IEvaluableComponent.js";
export default class SelectionOption implements IEvaluableComponent {
    component: IEvaluableComponent;
    weight: number;
    constructor(component: IEvaluableComponent, weight: number);
    toString(): string;
    evaluate(): string;
    evaluateAll(): string[];
}
export { SelectionOption };
