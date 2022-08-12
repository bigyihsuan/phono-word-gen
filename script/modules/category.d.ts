declare type Phoneme = string;
declare class Category {
    name: string;
    phonemes: Phoneme[];
    constructor(name: string, phonemes: Phoneme[]);
    toString(): string;
    isUnresolved(): boolean;
    add(other: Category): void;
    containedCategories(): string[];
    containedPhonemes(): Phoneme[];
}
declare type CategoryListing = Map<string, Category>;
declare function parseCategory(cat: string): Category;
declare function fillCategories(categories: CategoryListing): CategoryListing;
export { CategoryListing, Category, parseCategory, fillCategories, };
