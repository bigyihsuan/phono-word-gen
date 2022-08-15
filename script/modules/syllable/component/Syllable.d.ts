import { IEvaluableComponent } from "./IEvaluableComponent.js";
import IMatchable from "./IMatchable.js";
export default class Syllable implements IEvaluableComponent, IMatchable {
    components: IEvaluableComponent[];
    possibilities: string[];
    constructor(components: IEvaluableComponent[]);
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
    matches(word: string): boolean;
}
export { Syllable };
