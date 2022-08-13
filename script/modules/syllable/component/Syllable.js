export default class Syllable {
    components;
    constructor(components) {
        this.components = components;
    }
    evaluate() {
        return this.components.map((c) => c.evaluate()).join("");
    }
}
export { Syllable };
//# sourceMappingURL=Syllable.js.map