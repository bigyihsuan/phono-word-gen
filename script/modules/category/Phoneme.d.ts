export default class Phoneme {
    value: string;
    weight: number;
    constructor(phoneme: string);
    isCategoryName(): boolean;
    toString(): string;
}
export { Phoneme };
