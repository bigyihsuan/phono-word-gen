import { Phoneme } from "./Phoneme.js";
class Category {
    name;
    phonemes;
    weights = [];
    constructor(name, phonemes) {
        this.name = name;
        this.phonemes = phonemes;
        this.setWeights();
    }
    // see https://stackoverflow.com/a/55671924/8143168
    setWeights() {
        const unassignedSos = this.phonemes.filter((po) => po.weight < 0);
        const unassignedCount = unassignedSos.length;
        const totalWeight = this.phonemes
            .filter((po) => po.weight > 0) // only positive
            .map((po) => po.weight) // get the weights
            .reduce((p, w) => p + w, 0.0); // sum them
        const unassignedWeight = (1 - totalWeight) / unassignedCount;
        this.phonemes = this.phonemes.map((po) => {
            const s = po;
            s.weight = po.weight < 0 ? unassignedWeight : po.weight;
            return s;
        });
        for (let i = 0; i < this.phonemes.length; i += 1) {
            this.weights[i] = this.phonemes[i].weight + (this.weights[i - 1] || 0);
        }
    }
    // see https://stackoverflow.com/a/55671924/8143168
    getRandomChoice() {
        let i;
        const random = Math.random() * this.weights[this.weights.length - 1];
        for (i = 0; i < this.weights.length; i += 1) {
            if (this.weights[i] >= random) {
                break;
            }
        }
        return this.phonemes[i].value;
    }
    toString() {
        return `{${this.name}: [${this.phonemes.toString()}]}`;
    }
    isUnresolved() {
        return this.containedCategories().length > 0;
    }
    // add another category's phonemes to this one
    add(other) {
        this.phonemes = [...new Set([...this.containedPhonemes(), ...other.phonemes])];
    }
    // gets the names of the categories contained within this one
    containedCategories() {
        return this.phonemes.filter((p) => p.isCategoryName()).map((n) => n.value.substring(1));
    }
    containedPhonemes() {
        return this.phonemes.filter((p) => !p.isCategoryName());
    }
}
function parseCategory(cat) {
    let name = "";
    let phonemes = [];
    const split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides
    [name, phonemes] = [split[0], split[1].split(" ").map((s) => new Phoneme(s))];
    return new Category(name, phonemes);
}
// replaces all references to categories in all categories to their phonemes
function fillCategories(categories) {
    const filled = new Map();
    categories.forEach((cat, key) => {
        const newCat = new Category(cat.name, cat.phonemes);
        const catStack = [];
        while (newCat.isUnresolved()) {
            for (let i = 0; i < newCat.containedCategories().length; i += 1) {
                const catName = newCat.containedCategories()[i];
                if (catStack.includes(newCat.name)) {
                    // we've already seen ourselves, this is recursive
                    // recursion not allowed
                    throw new Error(`recusion is not allowed in categories:\n    ${newCat.name} and ${catName} contain each other`);
                }
                const c = categories.get(catName);
                if (c === undefined) {
                    throw new Error(`${catName} doesn't exist`);
                }
                newCat.add(c);
                catStack.push(catName);
            }
        }
        newCat.setWeights();
        filled.set(key, newCat);
    });
    return filled;
}
export { Category, parseCategory, fillCategories, };
//# sourceMappingURL=Category.js.map