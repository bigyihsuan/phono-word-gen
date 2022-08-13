import { Phoneme } from "./Phoneme.js";
import RandomlyChoosable from "../syllable/component/RandomlyChoosable.js";
declare class Category implements RandomlyChoosable {
    name: string;
    phonemes: Phoneme[];
    weights: number[];
    constructor(name: string, phonemes: Phoneme[]);
    setWeights(): void;
    getRandomChoice(): string;
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
