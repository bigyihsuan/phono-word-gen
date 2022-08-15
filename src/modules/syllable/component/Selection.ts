import { IEvaluableComponent } from "./IEvaluableComponent.js";
import { IRandomlyChoosable } from "./IRandomlyChoosable.js";
import { SelectionOption } from "./SelectionOption.js";

export default class Selection implements IRandomlyChoosable, IEvaluableComponent {
    options: SelectionOption[];

    weights: number[] = [];

    constructor(options: SelectionOption[]) {
        this.options = options;
        this.generateWeights();
    }

    // see https://stackoverflow.com/a/55671924/8143168
    generateWeights() {
        const unassignedSos = this.options.filter((so) => so.weight < 0);
        const unassignedCount = unassignedSos.length;
        const totalWeight = this.options
            .filter((so) => so.weight > 0) // only positive
            .map((so) => so.weight) // get the weights
            .reduce((p, w) => p + w, 0.0); // sum them
        const unassignedWeight = (1 - totalWeight) / unassignedCount;
        this.options = this.options.map((so) => {
            const s = so;
            s.weight = so.weight < 0 ? unassignedWeight : so.weight;
            return s;
        });
        for (let i = 0; i < this.options.length; i += 1) {
            this.weights[i] = this.options[i].weight + (this.weights[i - 1] || 0);
        }
    }

    // see https://stackoverflow.com/a/55671924/8143168
    getRandomChoice(): string {
        let i: number;
        const random = Math.random() * this.weights[this.weights.length - 1];
        for (i = 0; i < this.weights.length; i += 1) {
            if (this.weights[i] > random) {
                break;
            }
        }
        return this.options[i].component.evaluate();
    }

    evaluate(): string {
        return this.getRandomChoice();
    }

    evaluateAll(): string[] {
        return this.options.flatMap((o) => o.evaluateAll());
    }

    toString(): string {
        return `[${this.options.join(",")}]`;
    }
}

export { Selection };
