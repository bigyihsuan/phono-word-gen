export default class CategoryNode {
    category;
    constructor(category) {
        this.category = category;
    }
    getRandomChoice() {
        const randomIndex = Math.floor(Math.random() * this.category.phonemes.length);
        return this.category.phonemes[randomIndex];
    }
    evaluate() {
        return this.getRandomChoice();
    }
}
export { CategoryNode };
//# sourceMappingURL=CategoryNode.js.map