import { Category, CategoryListing, parseCategory } from "./modules/category.js";
import { tokenizeSyllable } from "./modules/syllable/lexer.js";
import { Token } from "./modules/syllable/token.js";
import { ParseError, parseSyllable, Syllable } from "./modules/syllable/parser.js";

const phonology = document.getElementById("phonology") as HTMLInputElement;
const submit = document.getElementById("submit") as HTMLButtonElement;
const wordOutput = document.getElementById("output");

const categories: CategoryListing = new Map<string, Category>();
let tokens: Token[];
let syllable: Syllable | ParseError;

submit?.addEventListener("click", () => {
    const lines = phonology?.value.replaceAll(/\n+/g, "\n").split("\n").filter((s) => s.length > 0);

    lines.forEach((l) => {
        const line = l.trim();
        if (line.match(/=/)) {
            const cat = parseCategory(line);
            categories.set(cat.name, cat);
        } else if (line.match(/syllable:/)) {
            // TODO: parse syllable structure
            tokens = tokenizeSyllable(line);
            syllable = parseSyllable(tokens.slice(), categories);
        }
    });

    wordOutput?.replaceChildren(); // clear the output for each run

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
        const p: HTMLParagraphElement = document.createElement("p");
        p.innerHTML = syllable.reason
        wordOutput?.append(p);
    } else {
        // console.log(syllable)
        for (let i = 0; i < 10; i += 1) {
            const p: HTMLParagraphElement = document.createElement("p");
            p.innerHTML = syllable.evaluate()
            wordOutput?.append(p);
        }
    }

});
