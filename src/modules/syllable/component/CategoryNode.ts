import { Category } from "../../category/Category.js";
import { IEvaluableComponent } from "./IEvaluableComponent.js";

export default class CategoryNode implements IEvaluableComponent {
    category: Category;

    constructor(category: Category) {
        this.category = category;
    }

    evaluate(): string {
        return this.category.getRandomChoice();
    }

    evaluateAll(): string[] {
        return this.category.evaluateAll();
    }

    toString(): string {
        return this.category.toString();
    }

    toRegex(): RegExp {
        return this.category.toRegex();
    }
}

export { CategoryNode };
