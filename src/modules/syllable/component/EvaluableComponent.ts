// represents a syllable component (category, optional, selection, raw phoneme, etc)
export default interface EvaluableComponent {
    // turn this into a string of phonemes for output
    evaluate(): string;
}

export { EvaluableComponent };
