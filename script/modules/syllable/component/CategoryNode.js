export default class CategoryNode {
    category;
    constructor(category) {
        this.category = category;
    }
    evaluate() {
        return this.category.getRandomChoice();
    }
    toString() {
        return this.category.toString();
    }
}
export { CategoryNode };
//# sourceMappingURL=CategoryNode.js.map