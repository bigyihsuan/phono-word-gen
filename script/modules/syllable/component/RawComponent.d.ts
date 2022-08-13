import { EvaluableComponent } from "./EvaluableComponent.js";
export default class RawComponent implements EvaluableComponent {
    component: string;
    constructor(component: string);
    evaluate(): string;
    toString(): string;
}
export { RawComponent };
