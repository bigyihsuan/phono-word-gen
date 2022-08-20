import { IEvaluableComponent } from "./IEvaluableComponent.js";
import IMatchable from "./IMatchable.js";

export default class Syllable implements IEvaluableComponent, IMatchable {
    components: IEvaluableComponent[];

    constructor(components: IEvaluableComponent[]) {
        this.components = components;
    }

    evaluate(): string {
        return this.components.map((c) => c.evaluate()).join("");
    }

    evaluateAll(): string[] {
        // let possibles: string[][] = [];
        const cartesian = (...a: any[]) => a.reduce(
            (l, b) => l.flatMap((d: any) => b.map((e: any) => [d, e].flat(2))),
        );
        // possibles = cartesian(this.components.map((p) => p.evaluateAll()));
        // return possibles.map((l) => l.join(""));
        const stuff = this.components.map((c) => c.evaluateAll());
        if (stuff.length === 0) {
            return [];
        }
        if (stuff.length === 1) {
            return stuff[0];
        }
        return cartesian(...stuff)
            .map(
                (t: string[]) => t.join(""),
            );
    }

    toString(): string {
        return `<${this.components.join(",")}>`;
    }

    matches(word: string): boolean {
        // return this.possibilities.some((pos) => word.match(pos));
        return this.toRegex().test(word);
    }

    toRegex(): RegExp {
        return new RegExp(this.components.map((c) => c.toRegex().source).join(""));
    }
}

export { Syllable };
