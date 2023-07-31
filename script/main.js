/* eslint-disable no-console */
import { fillCategory, parseCategory, } from "./modules/category/Category.js";
import tokenizeSyllable from "./modules/syllable/lexer.js";
import { ParseError } from "./modules/syllable/ParseError.js";
import { parseSyllable } from "./modules/syllable/parser.js";
import { Reject } from "./modules/postprocess/Reject.js";
import Replacement from "./modules/postprocess/Replacement.js";
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
const forceWordLimitElement = document.getElementById("forceWordLimit");
const duplicateAlertElement = document.getElementById("duplicateAlert");
const rejectedAlertElement = document.getElementById("rejectedAlert");
const copyOutputButton = document.getElementById("copyOutput");
let words = [];
let unsortedWords = [];
let outWords = words.map((syls) => syls.join(separateSyllablesElement.checked ? "." : ""));
let letters = [];
const rejects = [];
const replStrs = [];
let rejectRegexp;
let categories = new Map();
let syllable;
const replacements = [];
let replacedWords = 0;
separateSyllablesElement?.addEventListener("click", () => { main(true); });
submit?.addEventListener("click", () => { main(false); });
sortOutputElement?.addEventListener("change", (e) => {
    console.log("changed", { e });
    words = sortWords(letters);
    replacedWords = makeOutWords(replacedWords);
    renderOutput({
        debug: debugOutputElement.checked,
        separateSyllables: separateSyllablesElement.checked,
        letters,
        rejects,
        rejectRegexp,
        replacements,
        categories,
        syllable,
        words,
        outWords,
    });
});
copyOutputButton?.addEventListener("click", () => {
    wordOutputTextArea.select();
    wordOutputTextArea.setSelectionRange(0, wordOutputTextArea.value.length);
    navigator.clipboard.writeText(wordOutputTextArea.value);
});
main(false);
function main(keepPrevious) {
    let tokens = [];
    wordOutputTextArea.value = "";
    duplicateAlertElement.innerHTML = "";
    duplicateAlertElement.hidden = true;
    rejectedAlertElement.innerHTML = "";
    rejectedAlertElement.hidden = true;
    const wordCount = Number.parseInt(wordCountElement.value, 10);
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
                .trim()
                .split(/\s*\|\s*/)
                .map((s) => s
                .trim()
                // extract rejections in brackets
                .replace(/{(.+)}/, (_, g1) => g1))
                .filter((s) => s.length > 0));
        }
        else if (line.match(/letters:/)) {
            letters = line.replaceAll("letters:", "").split(" ").filter((e) => e.length > 0);
        }
        else if (line.match(/replace:/)) {
            replStrs.push(line);
        }
    });
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
    let rejectComps = [];
    if (rejects.length > 0) {
        try {
            rejectComps = rejects.map((r) => new Reject(r, categories));
            rejectRegexp = new RegExp(rejectComps.map((r) => r.toRegex().source).join("|"));
        }
        catch (e) {
            wordOutputTextArea.value = e;
            console.error(e);
            return;
        }
    }
    else {
        rejectRegexp = /$^/;
    }
    try {
        replStrs.forEach((r) => {
            replacements.push(new Replacement(r, categories));
        });
    }
    catch (e) {
        wordOutputTextArea.value += e;
        console.error(e);
        return;
    }
    // replacements.forEach((r) => {
    //     console.log(r);
    // });
    // if (debugOutputElement.checked) {
    //     wordOutputTextArea.value += `letters: ${letters.join(",")}\n\n`;
    //     wordOutputTextArea.value += `rejections: ${rejects.join(",")}\n\n`;
    // eslint-disable-next-line max-len
    //     wordOutputTextArea.value += `replacements:\n    ${replacements.map((r) => r.toString()).join("\n    ")}\n\n`;
    //     wordOutputTextArea.value += `reject regex: ${rejectRegexp}\n\n`;
    // eslint-disable-next-line max-len
    //     wordOutputTextArea.value += `categories: ${Array.from(categories).map((cn) => cn[1]).join("\n")}\n\n`;
    // }
    const sylLine = lines.find((l) => l.trim().match(/syllable:/))?.replaceAll("syllable:", "").trim();
    if (sylLine !== undefined) {
        tokens = tokenizeSyllable(sylLine);
        const maybeSyllable = parseSyllable(tokens.slice(), categories, sylLine);
        // if (debugOutputElement.checked) {
        //     wordOutputTextArea.value += `syllable: ${syllable}`;
        //     wordOutputTextArea.value += "\n---------------\n";
        // }
        if (maybeSyllable instanceof ParseError) {
            wordOutputTextArea.value += maybeSyllable.toString();
            console.error(maybeSyllable);
            return;
        }
        syllable = maybeSyllable;
        const possibleSyllableCount = syllable.evaluateAll().length;
        words = keepPrevious ? words : [];
        let rejectedCount = 0;
        let duplicateCount = 0;
        let generatedWords = 0;
        replacedWords = 0;
        while (words.length < wordCount) {
            const syls = generateWord(syllable, minSylCount, maxSylCount);
            generatedWords += 1;
            if (rejectRegexp.test(syls.join(""))) {
                // rejections
                rejectedCount += 1;
            }
            else if (!allowDuplicatesElement.checked && [...new Set(words.map((s) => s.join("")))].includes(syls.join(""))) {
                // duplicates
                duplicateCount += 1;
            }
            else {
                words.push(syls);
            }
            if (!forceWordLimitElement.checked && generatedWords >= wordCount) {
                break;
            }
            if (forceWordLimitElement.checked) {
                if (possibleSyllableCount * maxSylCount * maxSylCount <= wordCount
                    && generatedWords === possibleSyllableCount * maxSylCount * maxSylCount) {
                    rejectedAlertElement.innerHTML += `not enough possibilities: can only generate ${possibleSyllableCount * maxSylCount * maxSylCount}/${wordCount} desired words\n`;
                    rejectedAlertElement.hidden = false;
                    break;
                }
            }
        }
        rejectedAlertElement.innerHTML += `generated ${generatedWords} words`;
        rejectedAlertElement.hidden = false;
        if (rejectedCount > 0) {
            rejectedAlertElement.innerHTML += `, rejected ${rejectedCount} words`;
            rejectedAlertElement.hidden = false;
        }
        if (duplicateCount > 0) {
            duplicateAlertElement.innerHTML += `removed ${duplicateCount} duplicates`;
            duplicateAlertElement.hidden = false;
        }
        unsortedWords = words.slice();
        words = sortWords(letters);
        replacedWords = makeOutWords(replacedWords);
        rejectedAlertElement.innerHTML += `, replaced ${replacedWords} words`;
        rejectedAlertElement.hidden = false;
        // wordOutputTextArea.value += outWords.join("\n");
        renderOutput({
            debug: debugOutputElement.checked,
            separateSyllables: separateSyllablesElement.checked,
            letters,
            rejects,
            rejectRegexp,
            replacements,
            categories,
            syllable,
            words,
            outWords,
        });
    }
}
function renderOutput(data) {
    console.log("renderOutput", { data });
    const textArea = document.getElementById("outputText");
    textArea.value = data.outWords.join("\n");
    if (data.debug) {
        textArea.value += `letters: ${data.letters.join(",")}\n\n`;
        textArea.value += `rejections: ${data.rejects.join(",")}\n\n`;
        textArea.value += `reject regex: ${data.rejectRegexp}\n\n`;
        textArea.value += `categories: ${Array.from(data.categories).map((cn) => cn[1]).join("\n")}\n\n`;
        textArea.value += `syllable: ${data.syllable}`;
        textArea.value += "\n---------------\n";
    }
}
// generate a word as its syllables
function generateWord(syl, minSyllables, maxSyllables) {
    const outWord = [];
    const numSyllables = Math.max(minSyllables, Math.floor(maxSyllables - Math.random() * maxSyllables) + 1);
    for (let i = 0; i < numSyllables; i += 1) {
        outWord.push(syl.evaluate());
    }
    return outWord;
}
function compareWordsLetterwise(ls) {
    // convert to indexes per letter
    return (left, right) => {
        const letterIndexer = toIndexArray(ls);
        const leftIndexes = letterIndexer(left.lets);
        const rightIndexes = letterIndexer(right.lets);
        const smallestLength = Math.min(leftIndexes.length, rightIndexes.length);
        // compare letter-by-letter
        for (let i = 0; i < smallestLength; i += 1) {
            if (leftIndexes[i] < rightIndexes[i]) {
                return -1;
            }
            if (leftIndexes[i] > rightIndexes[i]) {
                return 1;
            }
        }
        // handle when the smaller is a subset of the larger
        // smaller word comes first
        if (leftIndexes.length < rightIndexes.length) {
            return -1;
        }
        return 1;
    };
}
function toIndexArray(ls) {
    return (wordLetters) => wordLetters.map((l) => ls.indexOf(l));
}
// tokenize a word (as its syllables) into a list of contained letters
function letterizeWord(word, ls) {
    return word.flatMap((syl) => letterizeSyllable(syl, ls));
}
// tokenize a syllable into letters
function letterizeSyllable(syl, ls) {
    const letterRegexp = new RegExp(`(${ls.slice().sort((a, b) => b.length - a.length).join("|")})`, "u");
    return syl.split(letterRegexp).filter((s) => s.length > 0);
}
function sortWords(ls) {
    console.log("sortWords", { ls });
    if (sortOutputElement.checked && ls.length > 0) {
        // sort based on letters
        // letters can be of any length
        // tokenize the words into their letters
        const letterizedWords = words.map((w) => ({ word: w, lets: letterizeWord(w, ls) }));
        // sort based on these letters
        const compare = compareWordsLetterwise(ls);
        words = letterizedWords.slice().sort(compare).map((obj) => obj.word);
    }
    else if (sortOutputElement.checked) {
        words = words.slice().sort();
    }
    else {
        words = unsortedWords;
    }
    return words;
}
function makeOutWords(rw) {
    let rws = rw;
    outWords = words.map((syls) => syls.join(separateSyllablesElement.checked ? "." : ""));
    if (!allowDuplicatesElement.checked) {
        const wordset = [...new Set(outWords)];
        if (wordset.length < outWords.length) {
            duplicateAlertElement.innerHTML += `removed ${outWords.length - wordset.length} duplicates`;
            duplicateAlertElement.hidden = false;
        }
        outWords = wordset;
    }
    // apply replacements
    outWords = outWords.map((word) => {
        let w = word;
        let applied = false;
        replacements.forEach((r) => {
            const out = r.apply(w);
            w = out.result;
            if (out.couldApply) {
                applied = true;
            }
        });
        if (applied) {
            rws += 1;
        }
        return w;
    });
    // resort after doing replacements
    if (sortOutputElement.checked) {
        const letterizedWords = outWords.map((w) => ({ word: [w], lets: letterizeWord([w], letters) }));
        outWords = letterizedWords.slice().map((obj) => obj.word).map((sarr) => sarr.join("")).sort();
    }
    return rws;
}
//# sourceMappingURL=main.js.map