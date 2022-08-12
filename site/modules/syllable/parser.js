import { CategoryNode } from "./component/CategoryNode.js";
import { OptionalComponent } from "./component/OptionalComponent.js";
import { ParseError } from "./ParseError.js";
import { RawComponent } from "./component/RawComponent.js";
import { Selection } from "./component/Selection.js";
import { SelectionOption } from "./component/SelectionOption.js";
import { Syllable } from "./component/Syllable.js";
import { SyllableExpr } from "./component/SyllableExpr.js";
import { CategoryToken, CommaToken, LbracketToken, LparenToken, RawComponentToken, RbracketToken, RparenToken, StarToken, WeightToken, } from "./token.js";
function parseSyllable(tokens, categories, sylStr) {
    const components = [];
    while (tokens.length > 0) {
        // break out when ending a selection or option
        if (tokens.at(0) instanceof CommaToken || tokens.at(0) instanceof RparenToken) {
            break;
        }
        const comp = parseSyllableExpr(tokens, categories, sylStr);
        if (comp instanceof ParseError) {
            return comp;
        }
        components.push(comp);
    }
    return new Syllable(components);
}
function parseSyllableExpr(tokens, categories, sylStr) {
    let component;
    const tok = tokens.shift();
    if (tok instanceof RawComponentToken) {
        component = parseRawComponent(tok);
    }
    else if (tok instanceof CategoryToken) {
        const cat = parseCategory(tok, categories, sylStr);
        if (cat instanceof ParseError) {
            return cat;
        }
        component = cat;
    }
    else if (tok instanceof LparenToken) {
        const opt = parseOptionalComponent(tokens, categories, sylStr);
        if (opt instanceof ParseError) {
            return opt;
        }
        component = opt;
    }
    else if (tok instanceof LbracketToken) {
        const sel = parseSelection(tokens, categories, sylStr);
        if (sel instanceof ParseError) {
            return sel;
        }
        component = sel;
    }
    else {
        return new ParseError(`SyllableExpr got some invalid token '${tok}'`, tok, sylStr);
    }
    return new SyllableExpr(component);
}
function parseSelection(tokens, categories, sylStr) {
    const components = [];
    // lbracket component star weight comma syllable comma
    while (tokens.length > 0) {
        const component = parseSyllable(tokens, categories, sylStr);
        if (component instanceof ParseError) {
            return component;
        }
        let option;
        // for weight
        const starOrComma = tokens.at(0);
        if (starOrComma instanceof StarToken) {
            option = parseSelectionOptionWithWeight(tokens, component, sylStr);
            if (option instanceof ParseError) {
                return option;
            }
            components.push(option);
        }
        else if (!(starOrComma instanceof CommaToken)) {
            return new ParseError(`Selection expected comma after component, got '${starOrComma}' instead`, starOrComma, sylStr);
        }
        else {
            // end of no-weight selection option, temporarily set weight to -1
            components.push(new SelectionOption(component, -1));
        }
        // consume comma
        tokens.shift();
        // check for end of selection
        const rbracket = tokens.at(0);
        if (rbracket === undefined) {
            return new ParseError(`Selection expected ']' to end, ran out of tokens instead`, rbracket, sylStr);
        }
        if (!(rbracket instanceof RbracketToken)) {
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
function parseSelectionOptionWithWeight(tokens, component, sylStr) {
    const star = tokens.shift();
    if (!(star instanceof StarToken)) {
        return new ParseError(`SelectionOption expected '*' after component, got '${star}' instead`, star, sylStr);
    }
    const weightTok = tokens.shift();
    if (!(weightTok instanceof WeightToken)) {
        return new ParseError(`SelectionOption expected a weight after '*', got '${weightTok}' instead`, weightTok, sylStr);
    }
    const weight = Number.parseFloat(weightTok.lexeme);
    if (Number.isNaN(weight)) {
        return new ParseError(`SelectionOption weight is not valid, got '${weightTok}' instead`, weightTok, sylStr);
    }
    return new SelectionOption(component, weight);
}
function parseOptionalComponent(tokens, categories, sylStr) {
    const component = parseSyllableExpr(tokens, categories, sylStr);
    if (component instanceof ParseError) {
        return component;
    }
    const rparen = tokens.shift();
    if (!(rparen instanceof RparenToken)) {
        return new ParseError(`OptionalComponent expected right paren, got '${rparen}' instead`, rparen, sylStr);
    }
    // optional weight: star weight
    const star = tokens.at(0);
    if (star instanceof StarToken) {
        tokens.shift();
        // expect weight
        const weightTok = tokens.shift();
        if (!(weightTok instanceof WeightToken)) {
            return new ParseError(`OptionalComponent expected weight after star, got '${weightTok}' instead`, weightTok, sylStr);
        }
        const weight = Number.parseFloat(weightTok.lexeme);
        if (Number.isNaN(weight)) {
            return new ParseError(`OptionalComponent weight is not valid, got '${weightTok}' instead`, weightTok, sylStr);
        }
        return new OptionalComponent(component, weight);
    }
    return new OptionalComponent(component);
}
function parseCategory(token, categories, sylStr) {
    const name = token.lexeme.replace("$", "");
    const category = categories.get(name);
    if (category === undefined) {
        return new ParseError(`Category ${name} not found`, token, sylStr);
    }
    return new CategoryNode(category);
}
function parseRawComponent(token) {
    return new RawComponent(token.lexeme);
}
export { Syllable, parseSyllable, ParseError, };
//# sourceMappingURL=parser.js.map