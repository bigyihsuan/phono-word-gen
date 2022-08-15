import { IEvaluableComponent } from "./IEvaluableComponent.js";

export default class RawComponent implements IEvaluableComponent {
    component: string = "";

    constructor(component: string) {
        this.component = component;
    }

    evaluate(): string {
        return this.component;
    }

    evaluateAll(): string[] {
        return [this.component];
    }

    toString(): string {
        return this.component;
    }
}

export { RawComponent };
