import { CategoryNode } from "./component/CategoryNode.js";
import { OptionalComponent } from "./component/OptionalComponent.js";
import { ParseError } from "./ParseError.js";
import { RawComponent } from "./component/RawComponent.js";
import { Selection } from "./component/Selection.js";
import { SelectionOption } from "./component/SelectionOption.js";
import { Syllable } from "./component/Syllable.js";
import { SyllableExpr } from "./component/SyllableExpr.js";
import { CategoryToken, CommaToken, LbracketToken, LcurlyToken, LparenToken, RawComponentToken, RbracketToken, RcurlyToken, RparenToken, StarToken, WeightToken, } from "./token.js";
function parseSyllable(tokens, categories, sylStr) {
    const components = [];
    while (tokens.length > 0) {
        // break out when ending a selection, option, or weight
        const t = tokens.at(0);
        if (t instanceof CommaToken
            || t instanceof RparenToken
            || t instanceof RcurlyToken
            || t instanceof StarToken) {
            break;
        }
        else if (tokens.at(0) instanceof LcurlyToken) {
            tokens.shift();
            const syl = parseGroupedComponents(tokens, categories, sylStr);
            if (syl instanceof ParseError) {
                return syl.within("SyllableExpr-LcurlyToken");
            }
            components.push(...syl.components);
        }
        else {
            const comp = parseSyllableExpr(tokens, categories, sylStr);
            if (comp instanceof ParseError) {
                return comp.within("Syllable");
            }
            components.push(comp);
        }
    }
    return new Syllable(components);
}
function parseGroupedComponents(tokens, categories, sylStr) {
    const component = parseSyllable(tokens, categories, sylStr);
    if (component instanceof ParseError) {
        return component.within("GroupedComponents");
    }
    const rcurly = tokens.shift();
    if (!(rcurly instanceof RcurlyToken)) {
        return new ParseError(`expected right curly, got '${rcurly}' instead`, rcurly, sylStr, "parseGroupedComponents");
    }
    return component;
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
            return cat.within("SyllableExpr-CategoryToken");
        }
        component = cat;
    }
    else if (tok instanceof LparenToken) {
        const opt = parseOptionalComponent(tokens, categories, sylStr);
        if (opt instanceof ParseError) {
            return opt.within("SyllableExpr-LparenToken");
        }
        component = opt;
    }
    else if (tok instanceof LbracketToken) {
        const sel = parseSelection(tokens, categories, sylStr);
        if (sel instanceof ParseError) {
            return sel.within("SyllableExpr-LbracketTok");
        }
        component = sel;
    }
    else {
        return new ParseError(`got some invalid token '${tok}'`, tok, sylStr, "SyllableExpr");
    }
    return new SyllableExpr(component);
}
function parseSelection(tokens, categories, sylStr) {
    const components = [];
    while (tokens.length > 0) {
        const option = parseSelectionOption(tokens, categories, sylStr);
        if (option instanceof ParseError) {
            return option.within("Selection");
        }
        components.push(option);
        // consume comma
        tokens.shift();
        // check for end of selection
        const rbracket = tokens.at(0);
        if (rbracket === undefined) {
            return new ParseError(`expected ']' to end, ran out of tokens instead`, rbracket, sylStr, "Selection-EndOfSelection");
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
// parse a selection option
// if there is a star following this, get the weight as well
function parseSelectionOption(tokens, categories, sylStr) {
    const component = parseSyllable(tokens, categories, sylStr);
    if (component instanceof ParseError) {
        return component.within("SelectionOption-Component").appendMessage("\n(Did you forget the trailing comma?)");
    }
    let weight = -1;
    if (tokens.at(0) instanceof StarToken) {
        const w = parseWeight(tokens, sylStr);
        if (w instanceof ParseError) {
            return w.within("SelectionOption-Weight");
        }
        weight = w;
    }
    return new SelectionOption(component, weight);
}
function parseWeight(tokens, sylStr) {
    const star = tokens.shift();
    if (!(star instanceof StarToken)) {
        return new ParseError(`expected '*' after component, got '${star}' instead`, star, sylStr, "SelectionOption-Star");
    }
    const weightTok = tokens.shift();
    if (!(weightTok instanceof WeightToken)) {
        return new ParseError(`expected a weight after '*', got '${weightTok}' instead`, weightTok, sylStr, "SelectionOption-Weight");
    }
    const weight = Number.parseFloat(weightTok.lexeme);
    if (Number.isNaN(weight)) {
        return new ParseError(`weight is not valid, got '${weightTok}' instead`, weightTok, sylStr, "SelectionOption-WeightEval");
    }
    return weight;
}
function parseOptionalComponent(tokens, categories, sylStr) {
    const component = parseSyllable(tokens, categories, sylStr);
    if (component instanceof ParseError) {
        return component.within("OptionalComponent");
    }
    const rparen = tokens.shift();
    if (!(rparen instanceof RparenToken)) {
        return new ParseError(`expected right paren, got '${rparen}' instead`, rparen, sylStr, "OptionalComponent-End");
    }
    // optional weight: star weight
    const star = tokens.at(0);
    if (star instanceof StarToken) {
        tokens.shift();
        // expect weight
        const weightTok = tokens.shift();
        if (!(weightTok instanceof WeightToken)) {
            return new ParseError(`expected weight after star, got '${weightTok}' instead`, weightTok, sylStr, "OptionalComponent-Weight");
        }
        const weight = Number.parseFloat(weightTok.lexeme);
        if (Number.isNaN(weight)) {
            return new ParseError(`weight is not valid, got '${weightTok}' instead`, weightTok, sylStr, "OptionalComponent-WeightEval");
        }
        return new OptionalComponent(component, weight);
    }
    return new OptionalComponent(component);
}
function parseCategory(token, categories, sylStr) {
    const name = token.lexeme.replace("$", "");
    const category = categories.get(name);
    if (category === undefined) {
        return new ParseError(`${name} not found`, token, sylStr, "Category");
    }
    return new CategoryNode(category);
}
function parseRawComponent(token) {
    return new RawComponent(token.lexeme);
}
export { Syllable, parseSyllable, ParseError, };
//# sourceMappingURL=parser.js.map