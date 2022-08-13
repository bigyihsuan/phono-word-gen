import { Category } from "../../category/Category.js";
import { EvaluableComponent } from "./EvaluableComponent.js";

export default class CategoryNode implements EvaluableComponent {
    category: Category;

    constructor(category: Category) {
        this.category = category;
    }

    evaluate(): string {
        return this.category.getRandomChoice();
    }
}

export { CategoryNode };
