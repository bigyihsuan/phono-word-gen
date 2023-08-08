/* eslint-disable no-console */
import {
    Category, CategoryListing, fillCategory, parseCategory,
} from "./modules/category/Category.js";
import tokenizeSyllable from "./modules/syllable/lexer.js";
import { ParseError } from "./modules/syllable/ParseError.js";
import { Syllable, parseSyllable } from "./modules/syllable/parser.js";
import { Token } from "./modules/syllable/token.js";
import { Reject } from "./modules/postprocess/Reject.js";
import Replacement from "./modules/postprocess/Replacement.js";

const phonology = document.getElementById("phonology") as HTMLInputElement;
const submit = document.getElementById("submit") as HTMLButtonElement;
const minSylCountElement = document.getElementById("minSylCount") as HTMLInputElement;
const maxSylCountElement = document.getElementById("maxSylCount") as HTMLInputElement;
const wordCountElement = document.getElementById("wordCount") as HTMLInputElement;
const sentenceCountElement = document.getElementById("sentenceCount") as HTMLInputElement;
const wordOutputTextArea = document.getElementById("outputText") as HTMLInputElement;

const allowDuplicatesElement = document.getElementById("allowDuplicates") as HTMLInputElement;
const sortOutputElement = document.getElementById("sortOutput") as HTMLInputElement;
const debugOutputElement = document.getElementById("debugOutput") as HTMLInputElement;
const separateSyllablesElement = document.getElementById("separateSyllables") as HTMLInputElement;
const forceWordLimitElement = document.getElementById("forceWordLimit") as HTMLInputElement;

const duplicateAlertElement = document.getElementById("duplicateAlert") as HTMLElement;
const rejectedAlertElement = document.getElementById("rejectedAlert") as HTMLElement;

const copyOutputButton = document.getElementById("copyOutput") as HTMLButtonElement;

const generateSentencesElement = document.getElementById("generateSentences") as HTMLInputElement;

const wordCountInputDiv = document.getElementById("wordCountInput") as HTMLDivElement;
const sentenceCountInputDiv = document.getElementById("sentenceCountInput") as HTMLDivElement;
sentenceCountInputDiv.hidden = true;

let wordCount = Number.parseInt(wordCountElement.value, 10);
let minSylCount = Number.parseInt(minSylCountElement.value, 10);
let maxSylCount = Number.parseInt(maxSylCountElement.value, 10);

let words: string[][] = [];
let unsortedWords: string[][] = [];
let outWords = words.map((syls) => syls.join(separateSyllablesElement.checked ? "." : ""));
let letters: string[] = [];
const rejects: string[] = [];
const replStrs: string[] = [];
let rejectRegexp: RegExp;
let categories: CategoryListing = new Map<string, Category>();
let syllable: Syllable;
const replacements: Replacement[] = [];
let replacedWords = 0;

separateSyllablesElement?.addEventListener("click", () => { main(true); });
generateSentencesElement?.addEventListener("change", () => {
    if (generateSentencesElement.checked) {
        // disable some word-related inputs
        allowDuplicatesElement.disabled = true;
        forceWordLimitElement.disabled = true;
        sortOutputElement.disabled = true;
        wordCountInputDiv.hidden = true;
        sentenceCountInputDiv.hidden = false;
    } else {
        allowDuplicatesElement.disabled = false;
        forceWordLimitElement.disabled = false;
        sortOutputElement.disabled = false;
        wordCountInputDiv.hidden = false;
        sentenceCountInputDiv.hidden = true;
    }
});

submit?.addEventListener("click", () => {
    wordCount = Number.parseInt(wordCountElement.value, 10);
    minSylCount = Number.parseInt(minSylCountElement.value, 10);
    maxSylCount = Number.parseInt(maxSylCountElement.value, 10);

    const lines = parseInput();

    try {
        categories = initCategories();
        rejectRegexp = initRejects();
        replStrs.forEach((r) => {
            replacements.push(new Replacement(r, categories));
        });
    } catch (e: any) {
        wordOutputTextArea.value = e;
        return;
    }

    const sylLine = lines.find((l) => l.trim().match(/syllable:/))?.replaceAll("syllable:", "").trim();
    if (sylLine === undefined) { return; }

    const tokens: Token[] = tokenizeSyllable(sylLine);
    const maybeSyllable = parseSyllable(tokens.slice(), categories, sylLine);

    if (maybeSyllable instanceof ParseError) {
        wordOutputTextArea.value += maybeSyllable.toString();
        return;
    }
    syllable = maybeSyllable as Syllable;

    if (generateSentencesElement.checked) {
        // generate sentenceCount sentences, of some random number of words each
        const sentences: string[] = [];
        for (let i = 0; i < sentenceCountElement.valueAsNumber; i += 1) {
            sentences.push(generateSentence());
        }
        wordOutputTextArea.value = sentences.join(" ");
    } else {
        main(false);
    }
});

function generateSentence(): string {
    const wc = 1 + peakedPowerLaw(15, 5, 50);
    const sentenceWords: string[] = [];
    for (let w = 0; w < wc; w += 1) {
        let word = generateWord(syllable, minSylCount, maxSylCount).join("");
        if (w === 0) {
            word = word.charAt(0).toUpperCase() + word.substring(1);
        }
        sentenceWords.push(word);
    }
    return `${sentenceWords.join(" ")}.`;
}

sortOutputElement?.addEventListener("change", () => {
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

function main(keepPrevious: boolean) {
    let tokens: Token[] = [];

    wordOutputTextArea.value = "";
    duplicateAlertElement.innerHTML = "";
    duplicateAlertElement.hidden = true;
    rejectedAlertElement.innerHTML = "";
    rejectedAlertElement.hidden = true;

    wordCount = Number.parseInt(wordCountElement.value, 10);
    minSylCount = Number.parseInt(minSylCountElement.value, 10);
    maxSylCount = Number.parseInt(maxSylCountElement.value, 10);

    if (maxSylCount < minSylCount) {
        minSylCountElement.value = maxSylCount.toString();
        minSylCount = maxSylCount;
    } else if (minSylCount > maxSylCount) {
        maxSylCountElement.value = minSylCount.toString();
        maxSylCount = minSylCount;
    }

    const lines = parseInput();

    try {
        categories = initCategories();
        rejectRegexp = initRejects();
        replStrs.forEach((r) => {
            replacements.push(new Replacement(r, categories));
        });
    } catch (e: any) {
        wordOutputTextArea.value = e;
        return;
    }

    const sylLine = lines.find((l) => l.trim().match(/syllable:/))?.replaceAll("syllable:", "").trim();
    if (sylLine === undefined) { return; }

    tokens = tokenizeSyllable(sylLine);
    const maybeSyllable = parseSyllable(tokens.slice(), categories, sylLine);

    if (maybeSyllable instanceof ParseError) {
        wordOutputTextArea.value += maybeSyllable.toString();
        return;
    }
    syllable = maybeSyllable as Syllable;

    words = keepPrevious ? words : [];
    replacedWords = makeWords(wordCount, minSylCount, maxSylCount);

    unsortedWords = words.slice();
    words = sortWords(letters);
    replacedWords = makeOutWords(replacedWords);

    rejectedAlertElement.innerHTML += `, replaced ${replacedWords} words`;
    rejectedAlertElement.hidden = false;

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

function parseInput(): string[] {
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
        } else if (line.match(/reject:/)) {
            rejects.push(
                ...line.replaceAll("reject:", "")
                    .trim()
                    .split(/\s*\|\s*/)
                    .map((s) => s
                        .trim()
                        // extract rejections in brackets
                        .replace(/{(.+)}/, (_, g1) => g1))
                    .filter((s) => s.length > 0),
            );
        } else if (line.match(/letters:/)) {
            letters = line.replaceAll("letters:", "").split(" ").filter((e) => e.length > 0);
        } else if (line.match(/replace:/)) {
            replStrs.push(line);
        }
    });
    return lines;
}

function initCategories(): Map<string, Category> {
    const maybeCats: CategoryListing = new Map<string, Category>();
    Array.from(categories).forEach(((nameCat) => {
        const cat = fillCategory(nameCat[0], categories);
        cat.setWeights();
        maybeCats.set(nameCat[0], cat);
    }));
    return maybeCats;
}

function initRejects(): RegExp {
    let rejectComps: Reject[] = [];
    if (rejects.length > 0) {
        try {
            rejectComps = rejects.map((r) => new Reject(r, categories));
            return new RegExp(rejectComps.map((r) => r.toRegex().source).join("|"));
        } catch (e: any) {
            wordOutputTextArea.value = e;
            console.error(e);
            return /(?:)/;
        }
    } else {
        return /$^/;
    }
}

function makeWords(count: number, minSyls: number, maxSyls: number): number {
    let rejectedCount = 0;
    let duplicateCount = 0;
    let generatedWords = 0;
    replacedWords = 0;
    const possibleSyllableCount = syllable.evaluateAll().length;

    while (words.length < count) {
        const syls = generateWord(syllable, minSyls, maxSyls);
        if (rejectRegexp.test(syls.join(""))) {
            // rejections
            rejectedCount += 1;
        } else if (
            !(!allowDuplicatesElement.disabled && allowDuplicatesElement.checked)
            && [...new Set(words.map((s) => s.join("")))].includes(syls.join(""))) {
            // duplicates
            duplicateCount += 1;
        } else {
            words.push(syls);
        }

        if (!(!forceWordLimitElement.disabled && forceWordLimitElement.checked)
                && generatedWords >= count) {
            break;
        }

        if (!forceWordLimitElement.disabled && forceWordLimitElement.checked) {
            const maxCount = possibleSyllableCount * maxSyls * maxSyls;
            if (maxCount <= count && generatedWords === maxCount) {
                const str = `not enough possibilities: can only generate ${maxCount}/${count} desired words\n`;
                rejectedAlertElement.innerHTML += str;
                rejectedAlertElement.hidden = false;
                break;
            }
        }
        generatedWords += 1;
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

    return replacedWords;
}

type OutputData = {
    debug: boolean,
    separateSyllables: boolean,
    letters: string[],
    rejects: string[],
    rejectRegexp: RegExp,
    replacements: Replacement[],
    categories: CategoryListing,
    syllable: Syllable,
    words: string[][],
    outWords: string[],
};

function renderOutput(data: OutputData) {
    const textArea = document.getElementById("outputText") as HTMLInputElement;
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
function generateWord(syl: Syllable, minSyllables: number, maxSyllables: number): string[] {
    const outWord: string[] = [];
    const numSyllables = minSyllables + powerLaw(maxSyllables, 50);
    for (let i = 0; i < numSyllables; i += 1) {
        outWord.push(syl.evaluate());
    }
    return outWord;
}

type letterizedWord = { word: string[], lets: string[] };

function compareWordsLetterwise(ls: string[]):
    (left: letterizedWord, right: letterizedWord) => number {
    // convert to indexes per letter

    return (left: letterizedWord, right: letterizedWord) => {
        const letterIndexer = toIndexArray(ls);
        const leftIndexes = letterIndexer(left.lets);
        const rightIndexes = letterIndexer(right.lets);
        const smallestLength = Math.min(leftIndexes.length, rightIndexes.length);

        // compare letter-by-letter
        for (let i = 0; i < smallestLength; i += 1) {
            if (leftIndexes[i] < rightIndexes[i]) {
                return -1;
            } if (leftIndexes[i] > rightIndexes[i]) {
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

function toIndexArray(ls: string[]): (letters: string[]) => number[] {
    return (wordLetters: string[]): number[] => wordLetters.map((l) => ls.indexOf(l));
}

// tokenize a word (as its syllables) into a list of contained letters
function letterizeWord(word: string[], ls: string[]): string[] {
    return word.flatMap((syl) => letterizeSyllable(syl, ls));
}

// tokenize a syllable into letters
function letterizeSyllable(syl: string, ls: string[]): string[] {
    const letterRegexp = new RegExp(`(${ls.slice().sort((a, b) => b.length - a.length).join("|")})`, "u");
    return syl.split(letterRegexp).filter((s) => s.length > 0);
}

function sortWords(ls: string[]): string[][] {
    if (!sortOutputElement.disabled && sortOutputElement.checked && ls.length > 0) {
        // sort based on letters
        // letters can be of any length
        // tokenize the words into their letters
        const letterizedWords: letterizedWord[] = words.map(
            (w) => ({ word: w, lets: letterizeWord(w, ls) }),
        );
            // sort based on these letters
        const compare = compareWordsLetterwise(ls);
        words = letterizedWords.slice().sort(compare).map((obj) => obj.word);
    } else if (!sortOutputElement.disabled && sortOutputElement.checked) {
        words = words.slice().sort();
    } else {
        words = unsortedWords;
    }
    return words;
}

function makeOutWords(rw: number): number {
    let rws = rw;
    outWords = words.map((syls) => syls.join(separateSyllablesElement.checked ? "." : ""));
    if (!(!allowDuplicatesElement.disabled && allowDuplicatesElement.checked)) {
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
        const letterizedWords: letterizedWord[] = outWords.map(
            (w) => ({ word: [w], lets: letterizeWord([w], letters) }),
        );
        outWords = letterizedWords.slice().map((obj) => obj.word).map((sarr) => sarr.join("")).sort();
    }
    return rws;
}

// based on code by Mark Rosenfelder for gen
// https://www.zompist.com/gen.html
function peakedPowerLaw(max: number, mode: number, prob: number): number {
    if (Math.random() > 0.5) {
        return mode + powerLaw(max - mode, prob);
    }
    return mode + powerLaw(mode + 1, prob);
}

function powerLaw(max: number, percentage: number): number {
    for (let r = 0; ; r = (r + 1) % max) {
        if (randomPercentage() < percentage) {
            return r;
        }
    }
}

function randomPercentage(): number {
    return Math.floor(Math.random() * 101);
}
