type Phoneme = string;

function isCategoryName(p: Phoneme): boolean {
    return p.at(0) === "$";
}

class Category {
    name: string;

    phonemes: Phoneme[];

    constructor(name: string, phonemes: Phoneme[]) {
        this.name = name;
        this.phonemes = phonemes;
    }

    toString(): string {
        return `{${this.name}: [${this.phonemes.toString()}]}`;
    }

    isUnresolved(): boolean {
        return this.containedCategories().length > 0;
    }

    // add another category's phonemes to this one
    add(other: Category) {
        this.phonemes = [...this.containedPhonemes(), ...other.phonemes];
    }

    // gets the names of the categories contained within this one
    containedCategories(): string[] {
        return this.phonemes.filter(isCategoryName).map((n) => n.substring(1));
    }

    containedPhonemes(): Phoneme[] {
        return this.phonemes.filter((p) => !isCategoryName(p));
    }
}

type CategoryListing = Map<string, Category>;

function parseCategory(cat: string): Category {
    let name: string = "";
    let phonemes: Phoneme[] = [];
    const split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides

    [name, phonemes] = [split[0], split[1].split(" ")];

    return new Category(name, phonemes);
}

// replaces all references to categories in all categories to their phonemes
function fillCategories(categories: CategoryListing): CategoryListing {
    const filled: CategoryListing = new Map<string, Category>();

    categories.forEach((cat, key) => {
        const newCat = new Category(cat.name, cat.phonemes);
        while (newCat.isUnresolved()) {
            newCat.containedCategories().forEach((catName) => {
                newCat.add(categories.get(catName)!);
            });
        }
        filled.set(key, newCat);
    });

    return filled;
}

export {
    CategoryListing,
    Category,
    parseCategory,
    fillCategories,
};
