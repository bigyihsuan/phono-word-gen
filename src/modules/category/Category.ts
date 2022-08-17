import { Phoneme } from "./Phoneme.js";
import IRandomlyChoosable from "../syllable/component/IRandomlyChoosable.js";
import IEvaluableComponent from "../syllable/component/IEvaluableComponent.js";

class Category implements IRandomlyChoosable, IEvaluableComponent {
    name: string;

    phonemes: Phoneme[];

    weights: number[] = [];

    constructor(name: string, phonemes: Phoneme[]) {
        this.name = name;
        this.phonemes = phonemes;
    }

    // see https://stackoverflow.com/a/55671924/8143168
    setWeights() {
        this.weights = [];

        const unsetCount = this.phonemes.filter((p) => !p.isManuallyWeighted).length;
        const totalWeight = this.phonemes
            .filter((p) => !Number.isNaN(p.weight))
            .map((p) => p.weight)
            .reduce((p, c) => p + c, 0.0);
        const defaultWeight = Math.max((1 - totalWeight) / unsetCount, 0);

        for (let i = 0; i < this.phonemes.length; i += 1) {
            if (Number.isNaN(this.phonemes[i].weight) || !this.phonemes[i].isManuallyWeighted) {
                // set this phoneme's weight to the default one
                this.phonemes[i].weight = defaultWeight;
            }
            // actually set the weight
            this.weights[i] = this.phonemes[i].weight + (this.weights[i - 1] || 0.0);
        }
    }

    // see https://stackoverflow.com/a/55671924/8143168
    getRandomChoice(): string {
        let i: number;
        const random = Math.random() * this.weights[this.weights.length - 1];
        for (i = 0; i < this.weights.length; i += 1) {
            if (this.weights[i] >= random) {
                break;
            }
        }
        return this.phonemes[i].value;
    }

    isUnresolved(): boolean {
        return this.containedCategories().length > 0;
    }

    // add another category's phonemes to this one
    add(other: Category) {
        this.phonemes = [
            ...new Map(
                [...this.containedPhonemes(), ...other.containedPhonemes()]
                    .map((phoneme) => [phoneme.value, phoneme]),
            ).values()];
    }

    // gets the names of the categories contained within this one
    containedCategories(): string[] {
        return this.phonemes
            .filter((p) => p.isCategoryName())
            .map((p) => p.copy())
            .map((n) => n.value.substring(1));
    }

    containedPhonemes(): Phoneme[] {
        return this.phonemes
            .filter((p) => !p.isCategoryName())
            .map((p) => p.copy());
    }

    containsPhoneme(phoneme: string): boolean {
        return this.phonemes.some((p) => p.value === phoneme);
    }

    evaluate(): string {
        return this.getRandomChoice();
    }

    evaluateAll(): string[] {
        return this.phonemes.flatMap((p) => p.value);
    }

    toString(): string {
        return `{${this.name}: [${this.phonemes.toString()}]}`;
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
function fillCategory(catName: string, categories: CategoryListing): Category {
    const cat = categories.get(catName);
    if (cat === undefined) {
        throw new Error(`name ${catName} doesn't exist`);
    }
    const newCat = new Category(catName, cat.containedPhonemes());
    cat.containedCategories().forEach((n) => {
        let innerCat = categories.get(n);
        if (innerCat === undefined) {
            throw new Error(`name ${n} doesn't exist`);
        }
        if (innerCat.isUnresolved()) {
            innerCat = fillCategory(innerCat.name, categories);
        }
        newCat.add(fillCategory(innerCat.name, categories));
    });
    return newCat;
}

export {
    CategoryListing,
    Category,
    parseCategory,
    fillCategory,
};
