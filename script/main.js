import { fillCategory, parseCategory, } from "./modules/category/Category.js";
import tokenizeSyllable from "./modules/syllable/lexer.js";
import { ParseError } from "./modules/syllable/ParseError.js";
import { Syllable, parseSyllable } from "./modules/syllable/parser.js";
const phonology = document.getElementById("phonology");
const submit = document.getElementById("submit");
const minSylCountElement = document.getElementById("minSylCount");
const maxSylCountElement = document.getElementById("maxSylCount");
const wordCountElement = document.getElementById("wordCount");
const wordOutputTextArea = document.getElementById("outputText");
const allowDuplicatesElement = document.getElementById("allowDuplicates");
const sortOutputElement = document.getElementById("sortOutput");
const debugOutputElement = document.getElementById("debugOutput");
const separateSyllablesElement = document.getElementById("separateSyllables");
const duplicateAlertElement = document.getElementById("duplicateAlert");
const rejectedAlertElement = document.getElementById("rejectedAlert");
submit?.addEventListener("click", () => {
    wordOutputTextArea.value = "";
    duplicateAlertElement.innerHTML = "";
    duplicateAlertElement.hidden = true;
    rejectedAlertElement.innerHTML = "";
    rejectedAlertElement.hidden = true;
    const wordCount = Number.parseInt(wordCountElement.value, 10);
    let categories = new Map();
    let tokens = [];
    let syllable;
    const rejects = [];
    let minSylCount = Number.parseInt(minSylCountElement.value, 10);
    let maxSylCount = Number.parseInt(maxSylCountElement.value, 10);
    if (maxSylCount < minSylCount) {
        minSylCountElement.value = maxSylCount.toString();
        minSylCount = maxSylCount;
    }
    else if (minSylCount > maxSylCount) {
        maxSylCountElement.value = minSylCount.toString();
        maxSylCount = minSylCount;
    }
    const lines = phonology?.value
        .replaceAll(/\n+/g, "\n") // remove extraneous newlines
        .replaceAll(/#.*/g, "") // remove comments
        .split("\n")
        .filter((s) => s.length > 0);
    lines.forEach((l) => {
        const line = l.trim();
        if (line.match(/=/)) {
            const cat = parseCategory(line);
            categories.set(cat.name, cat);
        }
        else if (line.match(/reject:/)) {
            rejects.push(...line.replaceAll("reject:", "")
                // extract rejections in curly brackets
                .split(/[<>]/)
                .map((s) => s
                .trim()
                .replace(/{(.+)}/, (_, g1) => g1))
                .filter((s) => s.length > 0));
        }
    });
    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `rejections: ${rejects.join(",")}\n`;
    }
    const maybeCats = new Map();
    try {
        Array.from(categories).forEach(((nameCat) => {
            const cat = fillCategory(nameCat[0], categories);
            cat.setWeights();
            maybeCats.set(nameCat[0], cat);
        }));
    }
    catch (e) {
        wordOutputTextArea.value = e;
        return;
    }
    categories = maybeCats;
    const rejectComps = rejects.map((r) => parseSyllable(tokenizeSyllable(r), categories, r)).filter((r) => r instanceof Syllable);
    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `parsed rejects:\n${rejectComps.join("\n")}\n`;
    }
    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `categories: ${Array.from(categories).map((cn) => cn[1]).join("\n")}\n`;
    }
    const sylLine = lines.find((l) => l.trim().match(/syllable:/))?.replaceAll("syllable:", "").trim();
    if (sylLine !== undefined) {
        tokens = tokenizeSyllable(sylLine);
        syllable = parseSyllable(tokens.slice(), categories, sylLine);
        if (debugOutputElement.checked) {
            wordOutputTextArea.value += `syllable: ${syllable}`;
            wordOutputTextArea.value += "\n---------------\n";
        }
        if (syllable instanceof ParseError) {
            wordOutputTextArea.value += syllable.toString();
        }
        else {
            if (debugOutputElement.checked) {
                wordOutputTextArea.value += `possibles: ${syllable.possibilities}`;
                wordOutputTextArea.value += "\n---------------\n";
            }
            let words = [];
            for (let _ = 0; _ < wordCount; _ += 1) {
                const outWord = [];
                const numSyllables = Math.max(minSylCount, Math.floor(maxSylCount - Math.random() * maxSylCount) + 1);
                for (let i = 0; i < numSyllables; i += 1) {
                    outWord.push(syllable.evaluate());
                }
                words.push(outWord);
            }
            if (!allowDuplicatesElement.checked) {
                const wordset = [...new Set(words)];
                if (wordset.length < words.length) {
                    duplicateAlertElement.innerHTML += `removed ${words.length - wordset.length} duplicates`;
                    duplicateAlertElement.hidden = false;
                }
                words = wordset;
            }
            const withoutRejected = words.filter((w) => !rejectComps.some((r) => (r instanceof Syllable) && r.matches(w.join(""))));
            if (words.length > withoutRejected.length) {
                rejectedAlertElement.innerHTML += `rejected ${words.length - withoutRejected.length} words`;
                rejectedAlertElement.hidden = false;
                words = withoutRejected;
            }
            let outWords = words.map((syls) => syls.join(separateSyllablesElement.checked ? "." : ""));
            if (sortOutputElement.checked) {
                outWords = outWords.sort();
            }
            wordOutputTextArea.value += outWords.join("\n");
        }
    }
});
//# sourceMappingURL=main.js.map