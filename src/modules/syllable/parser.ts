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
    token: Token
    sylStr: string

    constructor(reason: string, token: Token, sylStr: string) {
        this.reason = reason;
        this.token = token
        this.sylStr = sylStr;
    }

    toString(): string {
        let locationstr = "^".padStart(this.token.startingIndex + 1, ' ') + '^'.padStart(this.token.endingIndex - this.token.startingIndex - 1, '-');
        return [this.reason, this.sylStr, locationstr].join('\n')
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
    weight: number;
    // default weight is 50/50
    constructor(component: EvaluableComponent, weight: number = 0.5) {
        this.component = component;
        this.weight = weight;
    }

    getRandomChoice(): string {
        return Math.random() < this.weight ? this.component.evaluate() : "";
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

function parseSyllable(tokens: Token[], categories: CategoryListing, sylStr: string): Syllable | ParseError {
    let components: SyllableExpr[] = []
    while (tokens.length > 0) {
        // break out when ending a selection or option
        if (tokens.at(0) instanceof CommaToken || tokens.at(0) instanceof RparenToken) {
            break;
        }
        let comp = parseSyllableExpr(tokens, categories, sylStr)
        if (comp instanceof ParseError) {
            return comp
        }
        components.push(comp)
    }
    return new Syllable(components)
}

function parseSyllableExpr(tokens: Token[], categories: CategoryListing, sylStr: string): SyllableExpr | ParseError {
    let component: EvaluableComponent;
    let tok = tokens.shift()
    if (tok instanceof RawComponentToken) {
        component = parseRawComponent(tok, sylStr)
    } else if (tok instanceof CategoryToken) {
        let cat = parseCategory(tok, categories, sylStr)
        if (cat instanceof ParseError) {
            return cat
        }
        component = cat
    } else if (tok instanceof LparenToken) {
        let opt = parseOptionalComponent(tokens, categories, sylStr)
        if (opt instanceof ParseError) {
            return opt
        }
        component = opt
    } else if (tok instanceof LbracketToken) {
        let sel = parseSelection(tokens, categories, sylStr)
        if (sel instanceof ParseError) {
            return sel
        }
        component = sel
    } else {
        return new ParseError(`SyllableExpr got some invalid token '${tok}'`, tok!, sylStr)
    }
    return new SyllableExpr(component);
}

function parseSelection(tokens: Token[], categories: CategoryListing, sylStr: string): Selection | ParseError {
    let components: SelectionOption[] = []

    // lbracket component star weight comma syllable comma
    while (true) {
        let component = parseSyllable(tokens, categories, sylStr);
        if (component instanceof ParseError) {
            return component;
        }
        let option: SelectionOption | ParseError;
        // for weight
        let starOrComma = tokens.at(0)
        if (starOrComma instanceof StarToken) {
            option = parseSelectionOptionWithWeight(tokens, component, sylStr);
            if (option instanceof ParseError) {
                return option;
            }
            components.push(option);
        } else if (!(starOrComma instanceof CommaToken)) {
            return new ParseError(`Selection expected comma after component, got '${starOrComma}' instead`, starOrComma!, sylStr);
        } else {
            // end of no-weight selection option, temporarily set weight to -1
            components.push(new SelectionOption(component, -1))
        }
        // consume comma
        tokens.shift();
        // check for end of selection
        let rbracket = tokens.at(0);
        if (rbracket === undefined) {
            return new ParseError(`Selection expected ']' to end, ran out of tokens instead`, rbracket!, sylStr);
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

function parseSelectionOptionWithWeight(tokens: Token[], component: EvaluableComponent, sylStr: string): SelectionOption | ParseError {
    let star = tokens.shift();
    if (!(star instanceof StarToken)) {
        return new ParseError(`SelectionOption expected '*' after component, got '${star}' instead`, star!, sylStr);
    }
    let weightTok = tokens.shift();
    if (!(weightTok instanceof WeightToken)) {
        return new ParseError(`SelectionOption expected a weight after '*', got '${weightTok}' instead`, weightTok!, sylStr);
    }
    let weight = Number.parseFloat(weightTok.lexeme);
    if (weight === NaN) {
        return new ParseError(`SelectionOption weight is not valid, got '${weightTok}' instead`, weightTok!, sylStr);
    }
    return new SelectionOption(component, weight);
}

// TODO: optional component with weight
function parseOptionalComponent(tokens: Token[], categories: CategoryListing, sylStr: string): OptionalComponent | ParseError {
    let component = parseSyllableExpr(tokens, categories, sylStr);
    if (component instanceof ParseError) {
        return component
    }
    let rparen = tokens.shift()
    if (!(rparen instanceof RparenToken)) {
        return new ParseError(`OptionalComponent expected right paren, got '${rparen}' instead`, rparen!, sylStr)
    }
    // optional weight: star weight
    let star = tokens.at(0)
    if (star instanceof StarToken) {
        tokens.shift()
        // expect weight
        let weightTok = tokens.shift();
        if (!(weightTok instanceof WeightToken)) {
            return new ParseError(`OptionalComponent expected weight after star, got '${weightTok}' instead`, weightTok!, sylStr)
        }
        let weight = Number.parseFloat(weightTok.lexeme)
        if (weight === NaN) {
            return new ParseError(`OptionalComponent weight is not valid, got '${weightTok}' instead`, weightTok!, sylStr)
        }
        return new OptionalComponent(component, weight)
    }
    return new OptionalComponent(component)
}

function parseCategory(token: CategoryToken, categories: CategoryListing, sylStr: string): CategoryNode | ParseError {
    const name = token.lexeme.replace("$", "");
    const category = categories.get(name);
    if (category === undefined) {
        return new ParseError(`Category ${name} not found`, token!, sylStr);
    }
    return new CategoryNode(category);
}

function parseRawComponent(token: RawComponentToken, _: string): RawComponent {
    return new RawComponent(token.lexeme);
}

export {
    Syllable,
    parseSyllable,
    ParseError,
};
