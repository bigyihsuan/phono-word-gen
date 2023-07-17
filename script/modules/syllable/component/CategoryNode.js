export default class CategoryNode {
    category;
    constructor(category) {
        this.category = category;
    }
    evaluate() {
        return this.category.getRandomChoice();
    }
    evaluateAll() {
        return this.category.evaluateAll();
    }
    toString() {
        return this.category.toString();
    }
    toRegex() {
        return this.category.toRegex();
    }
}
export { CategoryNode };
//# sourceMappingURL=CategoryNode.js.map