declare const phonology: HTMLInputElement;
declare const submit: HTMLButtonElement;
declare class Category {
    name: string;
    phonemes: string[];
    constructor(name: string, phonemes: string[]);
}
declare function parseCategory(cat: string): Category;
