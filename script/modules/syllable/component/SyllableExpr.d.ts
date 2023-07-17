import { IEvaluableComponent } from "./IEvaluableComponent.js";
export default class SyllableExpr implements IEvaluableComponent {
    component: IEvaluableComponent;
    constructor(component: IEvaluableComponent);
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
    toRegex(): RegExp;
}
export { SyllableExpr };
