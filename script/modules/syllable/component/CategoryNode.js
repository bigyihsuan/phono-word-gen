export default class CategoryNode {
    category;
    constructor(category) {
        this.category = category;
    }
    evaluate() {
        return this.category.getRandomChoice();
    }
}
export { CategoryNode };
//# sourceMappingURL=CategoryNode.js.map