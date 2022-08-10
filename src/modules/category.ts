type Phoneme = string;

export class Category {
    name: string;
    phonemes: Phoneme[];

    constructor(name: string, phonemes: Phoneme[]) {
        this.name = name;
        this.phonemes = phonemes;
    }

    toString(): string {
        return `{${this.name}: [${this.phonemes.toString()}]}`
    }
}

export function parseCategory(cat: string): Category {
    let name: string = "";
    let phonemes: Phoneme[] = [];
    let split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides

    name = split[0];
    phonemes = split[1].split(" ")

    return new Category(name, phonemes);
}