import { EvaluableComponent } from "./EvaluableComponent.js";

export default class RawComponent implements EvaluableComponent {
    component: string = "";

    constructor(component: string) {
        this.component = component;
    }

    evaluate(): string {
        return this.component;
    }

    toString(): string {
        return this.component;
    }
}

export { RawComponent };
