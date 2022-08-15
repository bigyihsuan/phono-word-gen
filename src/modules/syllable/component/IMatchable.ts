// an interface for matching a string to a component
export default interface IMatchable {
    // returns whether a string can be output by this component
    matches(other: string): boolean;
}

export { IMatchable };
