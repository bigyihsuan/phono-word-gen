import { Category, CategoryListing } from "./../category.js";
import { CategoryToken, CommaToken, LbracketToken, LparenToken, RawComponentToken, RbracketToken, RparenToken, StarToken, Token, WeightToken } from "./token.js";

// represents a syllable component (category, optional, selection, raw phoneme, etc)
interface EvaluableComponent {
    // turn this into a string of phonemes for output
    evaluate(): string
}

// a component that can can be randomly chosen
interface RandomlyChoosable {
    // get and return a random phoneme/series of phonemes
    getRandomChoice(): string
}

class ParseError {
    reason: string;

    constructor(reason: string) {
        this.reason = reason;
    }
}

class Syllable implements EvaluableComponent {
    components: SyllableExpr[]
    constructor(components: SyllableExpr[]) {
        this.components = components
    }

    evaluate(): string {
        return this.components.map((c) => c.evaluate()).join("");
    }
}

class SyllableExpr implements EvaluableComponent {
    component: EvaluableComponent;

    constructor(component: EvaluableComponent) {
        this.component = component;
    }

    evaluate(): string {
        return this.component.evaluate()
    }
}

class Selection implements RandomlyChoosable, EvaluableComponent {
    options: SelectionOption[];

    weights: number[] = [];

    constructor(options: SelectionOption[]) {
        this.options = options;
        this.generateWeights();
    }

    // see https://stackoverflow.com/a/55671924/8143168
    generateWeights() {
        let i: number;

        let unassignedSos = this.options.filter((so) => so.weight < 0);
        let unassignedCount = unassignedSos.length
        let totalWeight = this.options
            .filter((so) => so.weight > 0) // only positive
            .map((so) => so.weight) // get the weights
            .reduce((p, w) => p + w, 0.0); // sum them
        let unassignedWeight = (1 - totalWeight) / unassignedCount;
        this.options.forEach((so) => {
            so.weight = so.weight < 0 ? unassignedWeight : so.weight;
        })
        for (i = 0; i < this.options.length; i++) {
            this.weights[i] = this.options[i].weight + (this.weights[i - 1] || 0);
        }
    }

    // see https://stackoverflow.com/a/55671924/8143168
    getRandomChoice(): string {
        let i: number;
        const random = Math.random() * this.weights[this.weights.length - 1];
        for (i = 0; i < this.weights.length; i++) {
            if (this.weights[i] > random) {
                break;
            }
        }
        return this.options[i].component.evaluate();
    }

    evaluate(): string {
        return this.getRandomChoice();
    }
}

class SelectionOption {
    component: EvaluableComponent;

    weight: number;

    constructor(component: EvaluableComponent, weight: number) {
        this.component = component;
        this.weight = weight;
    }
}

class OptionalComponent implements EvaluableComponent, RandomlyChoosable {
    component: EvaluableComponent;

    constructor(component: EvaluableComponent) {
        this.component = component;
    }

    getRandomChoice(): string {
        return Math.random() < 0.5 ? this.component.evaluate() : "";
    }

    evaluate(): string {
        return this.getRandomChoice();
    }
}

class CategoryNode implements RandomlyChoosable, EvaluableComponent {
    category: Category;

    constructor(category: Category) {
        this.category = category;
    }

    getRandomChoice(): string {
        const randomIndex = Math.floor(Math.random() * this.category.phonemes.length);
        return this.category.phonemes[randomIndex];
    }

    evaluate(): string {
        return this.getRandomChoice();
    }
}

class RawComponent implements EvaluableComponent {
    component: string = "";

    constructor(component: string) {
        this.component = component;
    }

    evaluate(): string {
        return this.component;
    }
}

function parseSyllable(tokens: Token[], categories: CategoryListing): Syllable | ParseError {
    let components: SyllableExpr[] = []
    while (tokens.length > 0) {
        let comp = parseSyllableExpr(tokens, categories)
        if (comp instanceof ParseError) {
            return comp
        }
        components.push(comp)
    }
    return new Syllable(components)
}

function parseSyllableExpr(tokens: Token[], categories: CategoryListing): SyllableExpr | ParseError {
    let component: EvaluableComponent;
    let tok = tokens.shift()
    if (tok instanceof RawComponentToken) {
        component = parseRawComponent(tok)
    } else if (tok instanceof CategoryToken) {
        let cat = parseCategory(tok, categories)
        if (cat instanceof ParseError) {
            return cat
        }
        component = cat
    } else if (tok instanceof LparenToken) {
        let opt = parseOptionalComponent(tokens, categories)
        if (opt instanceof ParseError) {
            return opt
        }
        component = opt
    } else if (tok instanceof LbracketToken) {
        let sel = parseSelection(tokens, categories)
        if (sel instanceof ParseError) {
            return sel
        }
        component = sel
    } else {
        return new ParseError(`SyllableExpr got some invalid token '${tok}'`)
    }
    return new SyllableExpr(component);
}

function parseSelection(tokens: Token[], categories: CategoryListing): Selection | ParseError {
    let components: SelectionOption[] = []

    // lbracket component star weight comma
    while (true) {
        let component = parseSyllableExpr(tokens, categories);
        if (component instanceof ParseError) {
            return component;
        }
        let option: SelectionOption | ParseError;
        // for weight
        let starOrComma = tokens.at(0)
        if (starOrComma instanceof StarToken) {
            option = parseSelectionOptionWithWeight(tokens, component);
            if (option instanceof ParseError) {
                return option;
            }
            components.push(option);
        } else if (!(starOrComma instanceof CommaToken)) {
            return new ParseError(`Selection expected comma after component, got '${starOrComma}' instead`);
        } else {
            // end of no-weight selection option, temporarily set weight to -1
            components.push(new SelectionOption(component, -1))
        }
        // consume comma
        tokens.shift();
        // check for end of selection
        let rbracket = tokens.at(0);
        if (rbracket === undefined) {
            return new ParseError(`Selection expected ']' to end, ran out of tokens instead`);
        } else if (!(rbracket instanceof RbracketToken)) {
            continue;
        } else {
            // end of selection
            // discard the right bracket and exit loop
            tokens.shift();
            break
        }
    }
    return new Selection(components);
}

function parseSelectionOptionWithWeight(tokens: Token[], component: EvaluableComponent): SelectionOption | ParseError {
    let star = tokens.shift();
    if (!(star instanceof StarToken)) {
        return new ParseError(`SelectionOption expected '*' after component, got '${star}' instead`);
    }
    let weightTok = tokens.shift();
    if (!(weightTok instanceof WeightToken)) {
        return new ParseError(`SelectionOption expected a weight after '*', got '${weightTok}' instead`);
    }
    let weight = Number.parseFloat(weightTok.lexeme);
    if (weight === NaN) {
        return new ParseError(`SelectionOption weight is not valid, got '${weightTok}' instead`);
    }
    return new SelectionOption(component, weight);
}

// TODO: optional component with weight
function parseOptionalComponent(tokens: Token[], categories: CategoryListing): OptionalComponent | ParseError {
    let component = parseSyllableExpr(tokens, categories);
    if (component instanceof ParseError) {
        return component
    }
    let rparen = tokens.shift()
    if (!(rparen instanceof RparenToken)) {
        return new ParseError(`OptionalComponent expected right paren, got '${rparen}' instead`)
    }
    return new OptionalComponent(component)
}

function parseCategory(token: CategoryToken, categories: CategoryListing): CategoryNode | ParseError {
    const name = token.lexeme.replace("$", "");
    const category = categories.get(name);
    if (category === undefined) {
        return new ParseError(`Category ${name} not found`);
    }
    return new CategoryNode(category);
}

function parseRawComponent(token: RawComponentToken): RawComponent {
    return new RawComponent(token.lexeme);
}

export {
    Syllable,
    parseSyllable,
    ParseError,
};
