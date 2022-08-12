import { EvaluableComponent } from "./EvaluableComponent.js";

export default class SelectionOption {
    component: EvaluableComponent;

    weight: number;

    constructor(component: EvaluableComponent, weight: number) {
        this.component = component;
        this.weight = weight;
    }
}

export { SelectionOption };
