export default class Phoneme {
    value;
    weight;
    isManuallyWeighted;
    constructor(phoneme) {
        // contains weight
        if (phoneme.match(/\*/)) {
            const p = phoneme.split(/\*/);
            const v = p.at(0);
            if (v !== undefined) {
                this.value = v;
            }
            else {
                throw new Error(`somehow bad phoneme ${phoneme}`);
            }
            const w = Number.parseFloat(p.at(-1));
            if (Number.isNaN(w)) {
                throw new Error(`invalid weight ${p.at(-1)}`);
            }
            this.weight = w;
            this.isManuallyWeighted = true;
        }
        else {
            this.value = phoneme;
            this.weight = Number.NaN; // all phonemes without weights are set to NaN initially
            this.isManuallyWeighted = false;
        }
    }
    isCategoryName() {
        return this.value.at(0) === "$";
    }
    toString() {
        return `${this.value}:${this.weight.toFixed(3)}`;
    }
    copy() {
        const p = new Phoneme(this.value);
        p.weight = this.weight;
        p.isManuallyWeighted = this.isManuallyWeighted;
        return p;
    }
}
export { Phoneme };
//# sourceMappingURL=Phoneme.js.map