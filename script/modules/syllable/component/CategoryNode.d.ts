import { Category } from "../../category/Category.js";
import { EvaluableComponent } from "./EvaluableComponent.js";
export default class CategoryNode implements EvaluableComponent {
    category: Category;
    constructor(category: Category);
    evaluate(): string;
}
export { CategoryNode };
