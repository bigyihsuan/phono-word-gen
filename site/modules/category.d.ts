declare type Phoneme = string;
export declare class Category {
    name: string;
    phonemes: Phoneme[];
    constructor(name: string, phonemes: Phoneme[]);
    toString(): string;
}
export declare function parseCategory(cat: string): Category;
export {};
