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
    toRegex() {
        return this.component.toRegex();
    }
    evaluate() {
        return this.component.evaluate();
    }
    evaluateAll() {
        return this.component.evaluateAll();
    }
}
export { SelectionOption };
//# sourceMappingURL=SelectionOption.js.map