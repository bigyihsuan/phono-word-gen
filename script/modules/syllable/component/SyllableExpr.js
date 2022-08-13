export default class SyllableExpr {
    component;
    constructor(component) {
        this.component = component;
    }
    evaluate() {
        return this.component.evaluate();
    }
    toString() {
        return this.component.toString();
    }
}
export { SyllableExpr };
//# sourceMappingURL=SyllableExpr.js.map