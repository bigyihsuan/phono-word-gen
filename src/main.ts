import { Category, CategoryListing, parseCategory } from "./modules/category.js";
import { tokenizeSyllable } from "./modules/syllable/lexer.js";
import { Token } from "./modules/syllable/token.js";
import { ParseError, parseSyllable, Syllable } from "./modules/syllable/parser.js";

const phonology = document.getElementById("phonology") as HTMLInputElement;
const submit = document.getElementById("submit") as HTMLButtonElement;
const wordOutput = document.getElementById("output");
const minSylCountElement = document.getElementById("minSylCount") as HTMLInputElement;
const maxSylCountElement = document.getElementById("maxSylCount") as HTMLInputElement;
const wordCountElement = document.getElementById("wordCount") as HTMLInputElement;
const wordOutputTextArea = document.getElementById("outputText") as HTMLInputElement;

const categories: CategoryListing = new Map<string, Category>();
let tokens: Token[];
let syllable: Syllable | ParseError;

submit?.addEventListener("click", () => {
    const lines = phonology?.value
        .replaceAll(/\n+/g, "\n") // remove extraneous newlines
        .replaceAll(/#.*/g, "") // remove comments
        .split("\n")
        .filter((s) => s.length > 0);
    let minSylCount = Number.parseInt(minSylCountElement.value);
    let maxSylCount = Number.parseInt(maxSylCountElement.value);
    if (minSylCount > maxSylCount) {
        minSylCountElement.value = maxSylCount.toString()
        minSylCount = maxSylCount
    }
    let wordCount = Number.parseInt(wordCountElement.value)

    lines.forEach((l) => {
        let line = l.trim();
        if (line.match(/=/)) {
            const cat = parseCategory(line);
            categories.set(cat.name, cat);
        } else if (line.match(/syllable:/)) {
            // TODO: parse syllable structure
            line = line.replaceAll("syllable:", "").trim();
            tokens = tokenizeSyllable(line);
            syllable = parseSyllable(tokens.slice(), categories, line);
        }
    });

    wordOutput?.replaceChildren(); // clear the output for each run
    wordOutputTextArea.value = "";

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
        wordOutputTextArea.value += syllable.toString();
    } else if (syllable !== undefined) {
        for (let _ = 0; _ < wordCount; _ += 1) {
            const p: HTMLParagraphElement = document.createElement("p");
            let outWord = ""
            let numSyllables = Math.max(minSylCount, Math.floor(maxSylCount - Math.random() * maxSylCount) + 1)
            for (let i = 0; i < numSyllables; i++) {
                outWord += syllable.evaluate()
            }
            p.innerHTML = outWord
            wordOutput?.append(p);
        }
    }

});
