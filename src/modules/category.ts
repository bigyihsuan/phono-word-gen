type Phoneme = string;

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
}

type CategoryListing = Map<string, Category>;

function parseCategory(cat: string): Category {
    let name: string = "";
    let phonemes: Phoneme[] = [];
    const split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides

    [name, phonemes] = [split[0], split[1].split(" ")];

    return new Category(name, phonemes);
}

export {
    CategoryListing,
    Category,
    parseCategory,
};
