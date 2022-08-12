import { EvaluableComponent } from "./EvaluableComponent.js";
import { RandomlyChoosable } from "./RandomlyChoosable.js";
import { SelectionOption } from "./SelectionOption.js";

export default class Selection implements RandomlyChoosable, EvaluableComponent {
    options: SelectionOption[];

    weights: number[] = [];

    constructor(options: SelectionOption[]) {
        this.options = options;
        this.generateWeights();
    }

    // see https://stackoverflow.com/a/55671924/8143168
    generateWeights() {
        let i: number;

        const unassignedSos = this.options.filter((so) => so.weight < 0);
        const unassignedCount = unassignedSos.length;
        const totalWeight = this.options
            .filter((so) => so.weight > 0) // only positive
            .map((so) => so.weight) // get the weights
            .reduce((p, w) => p + w, 0.0); // sum them
        const unassignedWeight = (1 - totalWeight) / unassignedCount;
        this.options.forEach((so) => {
            // eslint-disable-next-line no-param-reassign
            so.weight = so.weight < 0 ? unassignedWeight : so.weight;
        });
        for (i = 0; i < this.options.length; i += 1) {
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
}

export { Selection };
