import { Category } from "../../category/Category.js";
import { IEvaluableComponent } from "./IEvaluableComponent.js";
export default class CategoryNode implements IEvaluableComponent {
    category: Category;
    constructor(category: Category);
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
    toRegex(): RegExp;
}
export { CategoryNode };
