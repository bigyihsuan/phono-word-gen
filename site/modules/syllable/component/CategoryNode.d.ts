import { Category } from "../../category.js";
import { EvaluableComponent } from "./EvaluableComponent.js";
import { RandomlyChoosable } from "./RandomlyChoosable.js";
export default class CategoryNode implements RandomlyChoosable, EvaluableComponent {
    category: Category;
    constructor(category: Category);
    getRandomChoice(): string;
    evaluate(): string;
}
export { CategoryNode };
