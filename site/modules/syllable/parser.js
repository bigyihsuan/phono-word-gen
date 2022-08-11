import { CategoryToken, CommaToken, LbracketToken, LparenToken, RawComponentToken, RbracketToken, RparenToken, StarToken, WeightToken } from "./token.js";
class ParseError {
    constructor(reason) {
        this.reason = reason;
    }
}
class Syllable {
    constructor(components) {
        this.components = components;
    }
    evaluate() {
        return this.components.map((c) => c.evaluate()).join("");
    }
}
class SyllableExpr {
    constructor(component) {
        this.component = component;
    }
    evaluate() {
        return this.component.evaluate();
    }
}
class Selection {
    constructor(options) {
        this.weights = [];
        this.options = options;
        this.generateWeights();
    }
    // see https://stackoverflow.com/a/55671924/8143168
    generateWeights() {
        let i;
        let unassignedSos = this.options.filter((so) => so.weight < 0);
        let unassignedCount = unassignedSos.length;
        let totalWeight = this.options
            .filter((so) => so.weight > 0) // only positive
            .map((so) => so.weight) // get the weights
            .reduce((p, w) => p + w, 0.0); // sum them
        let unassignedWeight = (1 - totalWeight) / unassignedCount;
        this.options.forEach((so) => {
            so.weight = so.weight < 0 ? unassignedWeight : so.weight;
        });
        for (i = 0; i < this.options.length; i++) {
            this.weights[i] = this.options[i].weight + (this.weights[i - 1] || 0);
        }
    }
    // see https://stackoverflow.com/a/55671924/8143168
    getRandomChoice() {
        let i;
        const random = Math.random() * this.weights[this.weights.length - 1];
        for (i = 0; i < this.weights.length; i++) {
            if (this.weights[i] > random) {
                break;
            }
        }
        return this.options[i].component.evaluate();
    }
    evaluate() {
        return this.getRandomChoice();
    }
}
class SelectionOption {
    constructor(component, weight) {
        this.component = component;
        this.weight = weight;
    }
}
class OptionalComponent {
    constructor(component) {
        this.component = component;
    }
    getRandomChoice() {
        return Math.random() < 0.5 ? this.component.evaluate() : "";
    }
    evaluate() {
        return this.getRandomChoice();
    }
}
class CategoryNode {
    constructor(category) {
        this.category = category;
    }
    getRandomChoice() {
        const randomIndex = Math.floor(Math.random() * this.category.phonemes.length);
        return this.category.phonemes[randomIndex];
    }
    evaluate() {
        return this.getRandomChoice();
    }
}
class RawComponent {
    constructor(component) {
        this.component = "";
        this.component = component;
    }
    evaluate() {
        return this.component;
    }
}
function parseSyllable(tokens, categories) {
    let components = [];
    while (tokens.length > 0) {
        let comp = parseSyllableExpr(tokens, categories);
        if (comp instanceof ParseError) {
            return comp;
        }
        components.push(comp);
    }
    return new Syllable(components);
}
function parseSyllableExpr(tokens, categories) {
    let component;
    let tok = tokens.shift();
    if (tok instanceof RawComponentToken) {
        component = parseRawComponent(tok);
    }
    else if (tok instanceof CategoryToken) {
        let cat = parseCategory(tok, categories);
        if (cat instanceof ParseError) {
            return cat;
        }
        component = cat;
    }
    else if (tok instanceof LparenToken) {
        let opt = parseOptionalComponent(tokens, categories);
        if (opt instanceof ParseError) {
            return opt;
        }
        component = opt;
    }
    else if (tok instanceof LbracketToken) {
        let sel = parseSelection(tokens, categories);
        if (sel instanceof ParseError) {
            return sel;
        }
        component = sel;
    }
    else {
        return new ParseError(`SyllableExpr got some invalid token '${tok}'`);
    }
    return new SyllableExpr(component);
}
function parseSelection(tokens, categories) {
    let components = [];
    // lbracket component star weight comma
    while (true) {
        let component = parseSyllableExpr(tokens, categories);
        if (component instanceof ParseError) {
            return component;
        }
        let option;
        // for weight
        let starOrComma = tokens.at(0);
        if (starOrComma instanceof StarToken) {
            option = parseSelectionOptionWithWeight(tokens, component);
            if (option instanceof ParseError) {
                return option;
            }
            components.push(option);
        }
        else if (!(starOrComma instanceof CommaToken)) {
            return new ParseError(`Selection expected comma after component, got '${starOrComma}' instead`);
        }
        else {
            // end of no-weight selection option, temporarily set weight to -1
            components.push(new SelectionOption(component, -1));
        }
        // consume comma
        tokens.shift();
        // check for end of selection
        let rbracket = tokens.at(0);
        if (rbracket === undefined) {
            return new ParseError(`Selection expected ']' to end, ran out of tokens instead`);
        }
        else if (!(rbracket instanceof RbracketToken)) {
            continue;
        }
        else {
            // end of selection
            // discard the right bracket and exit loop
            tokens.shift();
            break;
        }
    }
    return new Selection(components);
}
function parseSelectionOptionWithWeight(tokens, component) {
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
function parseOptionalComponent(tokens, categories) {
    let component = parseSyllableExpr(tokens, categories);
    if (component instanceof ParseError) {
        return component;
    }
    let rparen = tokens.shift();
    if (!(rparen instanceof RparenToken)) {
        return new ParseError(`OptionalComponent expected right paren, got '${rparen}' instead`);
    }
    return new OptionalComponent(component);
}
function parseCategory(token, categories) {
    const name = token.lexeme.replace("$", "");
    const category = categories.get(name);
    if (category === undefined) {
        return new ParseError(`Category ${name} not found`);
    }
    return new CategoryNode(category);
}
function parseRawComponent(token) {
    return new RawComponent(token.lexeme);
}
export { Syllable, parseSyllable, ParseError, };
//# sourceMappingURL=parser.js.map