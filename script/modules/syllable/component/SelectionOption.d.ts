import { IEvaluableComponent } from "./IEvaluableComponent.js";
export default class SelectionOption implements IEvaluableComponent {
    component: IEvaluableComponent;
    weight: number;
    constructor(component: IEvaluableComponent, weight: number);
    toString(): string;
    toRegex(): RegExp;
    evaluate(): string;
    evaluateAll(): string[];
}
export { SelectionOption };
