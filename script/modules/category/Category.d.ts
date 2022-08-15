import { Phoneme } from "./Phoneme.js";
import IRandomlyChoosable from "../syllable/component/IRandomlyChoosable.js";
import IEvaluableComponent from "../syllable/component/IEvaluableComponent.js";
declare class Category implements IRandomlyChoosable, IEvaluableComponent {
    name: string;
    phonemes: Phoneme[];
    weights: number[];
    constructor(name: string, phonemes: Phoneme[]);
    setWeights(): void;
    getRandomChoice(): string;
    isUnresolved(): boolean;
    add(other: Category): void;
    containedCategories(): string[];
    containedPhonemes(): Phoneme[];
    containsPhoneme(phoneme: string): boolean;
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
}
declare type CategoryListing = Map<string, Category>;
declare function parseCategory(cat: string): Category;
declare function fillCategories(categories: CategoryListing): CategoryListing;
export { CategoryListing, Category, parseCategory, fillCategories, };
