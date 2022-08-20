export default class RawComponent {
    component = "";
    constructor(component) {
        this.component = component;
    }
    evaluate() {
        return this.component;
    }
    evaluateAll() {
        return [this.component];
    }
    toString() {
        return this.component;
    }
    toRegex() {
        return new RegExp(this.component);
    }
}
export { RawComponent };
//# sourceMappingURL=RawComponent.js.map