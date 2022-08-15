// represents a syllable component (category, optional, selection, raw phoneme, etc)
export default interface IEvaluableComponent {
    // turn this into a string of phonemes for output
    evaluate(): string;
    // create all possible strings
    evaluateAll(): string[];
    toString(): string;
}

export { IEvaluableComponent };
