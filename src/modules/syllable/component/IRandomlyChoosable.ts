// a component that can can be randomly chosen
export default interface IRandomlyChoosable {
    // get and return a random phoneme/series of phonemes
    getRandomChoice(): string;
}

export { IRandomlyChoosable };
