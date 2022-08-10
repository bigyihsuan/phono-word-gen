export class Category {
    constructor(name, phonemes) {
        this.name = name;
        this.phonemes = phonemes;
    }
    toString() {
        return `{${this.name}: [${this.phonemes.toString()}]}`;
    }
}
export function parseCategory(cat) {
    let name = "";
    let phonemes = [];
    let split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides
    name = split[0];
    phonemes = split[1].split(" ");
    return new Category(name, phonemes);
}
//# sourceMappingURL=category.js.map