export default class Phoneme {
    value: string;
    weight: number;
    isManuallyWeighted: boolean;
    constructor(phoneme: string);
    isCategoryName(): boolean;
    toString(): string;
    toRegex(): RegExp;
    copy(): Phoneme;
}
export { Phoneme };
