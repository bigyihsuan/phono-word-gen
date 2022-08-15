export default class SelectionOption {
    component;
    weight;
    constructor(component, weight) {
        this.component = component;
        this.weight = weight;
    }
    toString() {
        return `${this.component.toString()}:${this.weight}`;
    }
    evaluate() {
        return this.component.evaluate();
    }
}
export { SelectionOption };
//# sourceMappingURL=SelectionOption.js.map