export default class Phoneme {
    value;
    weight;
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
        }
        else {
            this.value = phoneme;
            this.weight = -1;
        }
    }
    isCategoryName() {
        return this.value.at(0) === "$";
    }
    toString() {
        return `${this.value}:${this.weight.toFixed(3)}`;
    }
}
export { Phoneme };
//# sourceMappingURL=Phoneme.js.map