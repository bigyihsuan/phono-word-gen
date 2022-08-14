import { EvaluableComponent } from "./EvaluableComponent.js";

export default class Syllable implements EvaluableComponent {
    components: EvaluableComponent[];

    constructor(components: EvaluableComponent[]) {
        this.components = components;
    }

    evaluate(): string {
        return this.components.map((c) => c.evaluate()).join("");
    }

    toString(): string {
        return `<${this.components.join(",")}>`;
    }
}

export { Syllable };
