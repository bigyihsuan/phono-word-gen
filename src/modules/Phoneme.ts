export default class Phoneme {
    value: string;

    weight: number;

    constructor(phoneme: string) {
        // contains weight
        if (phoneme.match(/\*/)) {
            const p = phoneme.split("*");
            this.value = p.at(0)!;
            const w = Number.parseFloat(p.at(-1)!);
            if (Number.isNaN(w)) {
                throw new Error(`invalid weight ${p.at(-1)!}`);
            }
            this.weight = w;
        } else {
            this.value = phoneme;
            this.weight = -1;
        }
    }

    isCategoryName(): boolean {
        return this.value.at(0) === "$";
    }
}

export { Phoneme };
