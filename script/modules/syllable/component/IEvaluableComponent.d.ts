export default interface IEvaluableComponent {
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
}
export { IEvaluableComponent };
