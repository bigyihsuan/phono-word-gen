function isCategoryName(p) {
    return p.at(0) === "$";
}
class Category {
    name;
    phonemes;
    constructor(name, phonemes) {
        this.name = name;
        this.phonemes = phonemes;
    }
    toString() {
        return `{${this.name}: [${this.phonemes.toString()}]}`;
    }
    isUnresolved() {
        return this.containedCategories().length > 0;
    }
    // add another category's phonemes to this one
    add(other) {
        this.phonemes = [...this.containedPhonemes(), ...other.phonemes];
    }
    // gets the names of the categories contained within this one
    containedCategories() {
        return this.phonemes.filter(isCategoryName).map((n) => n.substring(1));
    }
    containedPhonemes() {
        return this.phonemes.filter((p) => !isCategoryName(p));
    }
}
function parseCategory(cat) {
    let name = "";
    let phonemes = [];
    const split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides
    [name, phonemes] = [split[0], split[1].split(" ")];
    return new Category(name, phonemes);
}
// replaces all references to categories in all categories to their phonemes
function fillCategories(categories) {
    const filled = new Map();
    categories.forEach((cat, key) => {
        const newCat = new Category(cat.name, cat.phonemes);
        while (newCat.isUnresolved()) {
            newCat.containedCategories().forEach((catName) => {
                newCat.add(categories.get(catName));
            });
        }
        filled.set(key, newCat);
    });
    return filled;
}
export { Category, parseCategory, fillCategories, };
//# sourceMappingURL=category.js.map