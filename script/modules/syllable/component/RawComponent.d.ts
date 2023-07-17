import { IEvaluableComponent } from "./IEvaluableComponent.js";
export default class RawComponent implements IEvaluableComponent {
    component: string;
    constructor(component: string);
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
    toRegex(): RegExp;
}
export { RawComponent };
