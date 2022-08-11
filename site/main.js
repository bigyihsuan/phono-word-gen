import { parseCategory } from "./modules/category.js";
import { tokenizeSyllable } from "./modules/syllable/lexer.js";
import { ParseError, parseSyllable } from "./modules/syllable/parser.js";
const phonology = document.getElementById("phonology");
const submit = document.getElementById("submit");
const wordOutput = document.getElementById("output");
const minSylCountElement = document.getElementById("minSylCount");
const maxSylCountElement = document.getElementById("maxSylCount");
const wordCountElement = document.getElementById("wordCount");
const wordOutputTextArea = document.getElementById("outputText");
const categories = new Map();
let tokens;
let syllable;
submit === null || submit === void 0 ? void 0 : submit.addEventListener("click", () => {
    const lines = phonology === null || phonology === void 0 ? void 0 : phonology.value.replaceAll(/\n+/g, "\n").replaceAll(/#.*/g, "").split("\n").filter((s) => s.length > 0);
    let minSylCount = Number.parseInt(minSylCountElement.value);
    let maxSylCount = Number.parseInt(maxSylCountElement.value);
    if (minSylCount > maxSylCount) {
        minSylCountElement.value = maxSylCount.toString();
        minSylCount = maxSylCount;
    }
    let wordCount = Number.parseInt(wordCountElement.value);
    lines.forEach((l) => {
        let line = l.trim();
        if (line.match(/=/)) {
            const cat = parseCategory(line);
            categories.set(cat.name, cat);
        }
        else if (line.match(/syllable:/)) {
            // TODO: parse syllable structure
            line = line.replaceAll("syllable:", "").trim();
            tokens = tokenizeSyllable(line);
            syllable = parseSyllable(tokens.slice(), categories, line);
        }
    });
    wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.replaceChildren(); // clear the output for each run
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
    }
    else if (syllable !== undefined) {
        for (let _ = 0; _ < wordCount; _ += 1) {
            const p = document.createElement("p");
            let outWord = "";
            let numSyllables = Math.max(minSylCount, Math.floor(maxSylCount - Math.random() * maxSylCount) + 1);
            for (let i = 0; i < numSyllables; i++) {
                outWord += syllable.evaluate();
            }
            p.innerHTML = outWord;
            wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.append(p);
        }
    }
});
//# sourceMappingURL=main.js.map