export default class Syllable {
    components;
    constructor(components) {
        this.components = components;
    }
    evaluate() {
        return this.components.map((c) => c.evaluate()).join("");
    }
    evaluateAll() {
        // let possibles: string[][] = [];
        const cartesian = (...a) => a.reduce((l, b) => l.flatMap((d) => b.map((e) => [d, e].flat(2))));
        // possibles = cartesian(this.components.map((p) => p.evaluateAll()));
        // return possibles.map((l) => l.join(""));
        const stuff = this.components.map((c) => c.evaluateAll());
        if (stuff.length === 0) {
            return [];
        }
        if (stuff.length === 1) {
            return stuff[0];
        }
        return cartesian(...stuff)
            .map((t) => t.join(""));
    }
    toString() {
        return `<${this.components.join(",")}>`;
    }
    matches(word) {
        // return this.possibilities.some((pos) => word.match(pos));
        return this.toRegex().test(word);
    }
    toRegex() {
        return new RegExp(this.components.map((c) => c.toRegex().source).join(""));
    }
}
export { Syllable };
//# sourceMappingURL=Syllable.js.map