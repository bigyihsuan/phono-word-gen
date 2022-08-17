export default class Phoneme {
    value: string;

    weight: number;

    isManuallyWeighted: boolean;

    constructor(phoneme: string) {
        // contains weight
        if (phoneme.match(/\*/)) {
            const p = phoneme.split(/\*/);
            const v = p.at(0);
            if (v !== undefined) {
                this.value = v;
            } else {
                throw new Error(`somehow bad phoneme ${phoneme}`);
            }
            const w = Number.parseFloat(p.at(-1)!);
            if (Number.isNaN(w)) {
                throw new Error(`invalid weight ${p.at(-1)!}`);
            }
            this.weight = w;
            this.isManuallyWeighted = true;
        } else {
            this.value = phoneme;
            this.weight = Number.NaN; // all phonemes without weights are set to NaN initially
            this.isManuallyWeighted = false;
        }
    }

    isCategoryName(): boolean {
        return this.value.at(0) === "$";
    }

    toString(): string {
        return `${this.value}:${this.weight.toFixed(3)}`;
    }

    copy(): Phoneme {
        const p = new Phoneme(this.value);
        p.weight = this.weight;
        p.isManuallyWeighted = this.isManuallyWeighted;
        return p;
    }
}

export { Phoneme };
