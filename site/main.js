import { parseCategory } from "./modules/category.js";
import tokenizeSyllable from "./modules/syllable/lexer.js";
import { ParseError } from "./modules/syllable/ParseError.js";
import { parseSyllable } from "./modules/syllable/parser.js";
const phonology = document.getElementById("phonology");
const submit = document.getElementById("submit");
const wordOutput = document.getElementById("output");
const minSylCountElement = document.getElementById("minSylCount");
const maxSylCountElement = document.getElementById("maxSylCount");
const wordCountElement = document.getElementById("wordCount");
const wordOutputTextArea = document.getElementById("outputText");
const allowDuplicatesElement = document.getElementById("allowDuplicates");
const categories = new Map();
let tokens;
let syllable;
submit === null || submit === void 0 ? void 0 : submit.addEventListener("click", () => {
    const lines = phonology === null || phonology === void 0 ? void 0 : phonology.value.replaceAll(/\n+/g, "\n").replaceAll(/#.*/g, "").split("\n").filter((s) => s.length > 0);
    const minSylCount = Number.parseInt(minSylCountElement.value, 10);
    let maxSylCount = Number.parseInt(maxSylCountElement.value, 10);
    if (minSylCount > maxSylCount) {
        maxSylCountElement.value = minSylCount.toString();
        maxSylCount = minSylCount;
    }
    const wordCount = Number.parseInt(wordCountElement.value, 10);
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
    wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.append(document.createElement("hr"));
    if (syllable instanceof ParseError) {
        wordOutputTextArea.value += syllable.toString();
    }
    else if (syllable !== undefined) {
        let words = [];
        for (let _ = 0; _ < wordCount; _ += 1) {
            let outWord = "";
            const numSyllables = Math.max(minSylCount, Math.floor(maxSylCount - Math.random() * maxSylCount) + 1);
            for (let i = 0; i < numSyllables; i += 1) {
                outWord += syllable.evaluate();
            }
            words.push(outWord);
        }
        if (!allowDuplicatesElement.checked) {
            words = [...new Set(words)];
        }
        wordOutputTextArea.value = words.join("\n");
    }
});
//# sourceMappingURL=main.js.map