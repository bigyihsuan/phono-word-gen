export default interface IEvaluableComponent {
    evaluate(): string;
    evaluateAll(): string[];
    toString(): string;
    toRegex(): RegExp;
}
export { IEvaluableComponent };
