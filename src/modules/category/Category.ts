import { Phoneme } from "./Phoneme.js";
import RandomlyChoosable from "../syllable/component/RandomlyChoosable.js";

class Category implements RandomlyChoosable {
    name: string;

    phonemes: Phoneme[];

    weights: number[] = [];

    constructor(name: string, phonemes: Phoneme[]) {
        this.name = name;
        this.phonemes = phonemes;
        this.setWeights();
    }

    // see https://stackoverflow.com/a/55671924/8143168
    setWeights() {
        const unassignedSos = this.phonemes.filter((so) => so.weight < 0);
        const unassignedCount = unassignedSos.length;
        const totalWeight = this.phonemes
            .filter((so) => so.weight > 0) // only positive
            .map((so) => so.weight) // get the weights
            .reduce((p, w) => p + w, 0.0); // sum them
        const unassignedWeight = (1 - totalWeight) / unassignedCount;
        this.phonemes = this.phonemes.map((so) => {
            const s = so;
            s.weight = so.weight < 0 ? unassignedWeight : so.weight;
            return s;
        });
        for (let i = 0; i < this.phonemes.length; i += 1) {
            this.weights[i] = this.phonemes[i].weight + (this.weights[i - 1] || 0);
        }
    }

    // see https://stackoverflow.com/a/55671924/8143168
    getRandomChoice(): string {
        let i: number;
        const random = Math.random() * this.weights[this.weights.length - 1];
        for (i = 0; i < this.weights.length; i += 1) {
            if (this.weights[i] > random) {
                break;
            }
        }
        return this.phonemes[i].value;
    }

    toString(): string {
        return `{${this.name}: [${this.phonemes.toString()}]}`;
    }

    isUnresolved(): boolean {
        return this.containedCategories().length > 0;
    }

    // add another category's phonemes to this one
    add(other: Category) {
        this.phonemes = [...new Set([...this.containedPhonemes(), ...other.phonemes])];
    }

    // gets the names of the categories contained within this one
    containedCategories(): string[] {
        return this.phonemes.filter((p) => p.isCategoryName()).map((n) => n.value.substring(1));
    }

    containedPhonemes(): Phoneme[] {
        return this.phonemes.filter((p) => !p.isCategoryName());
    }
}

type CategoryListing = Map<string, Category>;

function parseCategory(cat: string): Category {
    let name: string = "";
    let phonemes: Phoneme[] = [];
    const split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides

    [name, phonemes] = [split[0], split[1].split(" ").map((s) => new Phoneme(s))];

    return new Category(name, phonemes);
}

// replaces all references to categories in all categories to their phonemes
function fillCategories(categories: CategoryListing): CategoryListing | Error {
    const filled: CategoryListing = new Map<string, Category>();

    try {
        categories.forEach((cat, key) => {
            const newCat = new Category(cat.name, cat.phonemes);
            const catStack: string[] = [];
            while (newCat.isUnresolved()) {
                for (let i = 0; i < newCat.containedCategories().length; i += 1) {
                    const catName = newCat.containedCategories()[i];
                    if (catStack.includes(newCat.name)) {
                        // we've already seen ourselves, this is recursive
                        // recursion not allowed
                        throw new Error(`recusion is not allowed in categories:\n    ${newCat.name} and ${catName} contain each other`);
                    }
                    newCat.add(categories.get(catName)!);
                    catStack.push(catName);
                }
            }
            filled.set(key, newCat);
        });
    } catch (e: any) {
        return e;
    }

    return filled;
}

export {
    CategoryListing,
    Category,
    parseCategory,
    fillCategories,
};
