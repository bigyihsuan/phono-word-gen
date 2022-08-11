declare type Phoneme = string;
declare class Category {
    name: string;
    phonemes: Phoneme[];
    constructor(name: string, phonemes: Phoneme[]);
    toString(): string;
}
declare type CategoryListing = Map<string, Category>;
declare function parseCategory(cat: string): Category;
export { CategoryListing, Category, parseCategory, };
