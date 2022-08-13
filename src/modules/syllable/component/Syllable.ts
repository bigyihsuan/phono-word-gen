import { EvaluableComponent } from "./EvaluableComponent.js";
import { SyllableExpr } from "./SyllableExpr.js";

export default class Syllable implements EvaluableComponent {
    components: SyllableExpr[];

    constructor(components: SyllableExpr[]) {
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
