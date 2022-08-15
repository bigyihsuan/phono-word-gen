export default class SyllableExpr {
    component;
    constructor(component) {
        this.component = component;
    }
    evaluate() {
        return this.component.evaluate();
    }
    evaluateAll() {
        return this.component.evaluateAll();
    }
    toString() {
        return this.component.toString();
    }
}
export { SyllableExpr };
//# sourceMappingURL=SyllableExpr.js.map