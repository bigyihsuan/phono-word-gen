export default class OptionalComponent {
    component;
    weight;
    // default weight is 50/50
    constructor(component, weight = 0.5) {
        this.component = component;
        this.weight = weight;
    }
    getRandomChoice() {
        return Math.random() < this.weight ? this.component.evaluate() : "";
    }
    evaluate() {
        return this.getRandomChoice();
    }
}
export { OptionalComponent };
//# sourceMappingURL=OptionalComponent.js.map