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
}
export { CategoryNode };
//# sourceMappingURL=CategoryNode.js.map