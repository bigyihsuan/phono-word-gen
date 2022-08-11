import { parseCategory } from "./modules/category.js";
import { tokenizeSyllable } from "./modules/syllable/lexer.js";
import { ParseError, parseSyllable } from "./modules/syllable/parser.js";
const phonology = document.getElementById("phonology");
const submit = document.getElementById("submit");
const wordOutput = document.getElementById("output");
const categories = new Map();
let tokens;
let syllable;
submit === null || submit === void 0 ? void 0 : submit.addEventListener("click", () => {
    const lines = phonology === null || phonology === void 0 ? void 0 : phonology.value.replaceAll(/\n+/g, "\n").split("\n").filter((s) => s.length > 0);
    lines.forEach((l) => {
        const line = l.trim();
        if (line.match(/=/)) {
            const cat = parseCategory(line);
            categories.set(cat.name, cat);
        }
        else if (line.match(/syllable:/)) {
            // TODO: parse syllable structure
            tokens = tokenizeSyllable(line);
            syllable = parseSyllable(tokens.slice(), categories);
        }
    });
    wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.replaceChildren(); // clear the output for each run
    // categories.forEach((c) => {
    //     const p: HTMLParagraphElement = document.createElement("p");
    //     p.innerHTML = c.toString();
    //     wordOutput?.append(p);
    // });
    // wordOutput?.append(document.createElement("hr"));
    // tokens.slice().forEach((t) => {
    //     const p: HTMLParagraphElement = document.createElement("p");
    //     p.innerHTML = t.toString();
    //     wordOutput?.append(p);
    // });
    // wordOutput?.append(document.createElement("hr"));
    if (syllable instanceof ParseError) {
        const p = document.createElement("p");
        p.innerHTML = syllable.reason;
        wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.append(p);
    }
    else {
        // console.log(syllable)
        for (let i = 0; i < 10; i += 1) {
            const p = document.createElement("p");
            p.innerHTML = syllable.evaluate();
            wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.append(p);
        }
    }
});
//# sourceMappingURL=main.js.map