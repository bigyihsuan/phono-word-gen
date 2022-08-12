import { Category } from "../../category.js";
import { EvaluableComponent } from "./EvaluableComponent.js";
import { RandomlyChoosable } from "./RandomlyChoosable.js";

export default class CategoryNode implements RandomlyChoosable, EvaluableComponent {
    category: Category;

    constructor(category: Category) {
        this.category = category;
    }

    getRandomChoice(): string {
        const randomIndex = Math.floor(Math.random() * this.category.phonemes.length);
        return this.category.phonemes[randomIndex];
    }

    evaluate(): string {
        return this.getRandomChoice();
    }
}

export { CategoryNode };
