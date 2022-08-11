class Category {
    constructor(name, phonemes) {
        this.name = name;
        this.phonemes = phonemes;
    }
    toString() {
        return `{${this.name}: [${this.phonemes.toString()}]}`;
    }
}
function parseCategory(cat) {
    let name = "";
    let phonemes = [];
    const split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides
    [name, phonemes] = [split[0], split[1].split(" ")];
    return new Category(name, phonemes);
}
export { Category, parseCategory, };
//# sourceMappingURL=category.js.map