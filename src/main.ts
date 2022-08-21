import {
    Category, CategoryListing, fillCategory, parseCategory,
} from "./modules/category/Category.js";
import tokenizeSyllable from "./modules/syllable/lexer.js";
import { ParseError } from "./modules/syllable/ParseError.js";
import { Syllable, parseSyllable } from "./modules/syllable/parser.js";
import { Token } from "./modules/syllable/token.js";

const phonology = document.getElementById("phonology") as HTMLInputElement;
const submit = document.getElementById("submit") as HTMLButtonElement;
const minSylCountElement = document.getElementById("minSylCount") as HTMLInputElement;
const maxSylCountElement = document.getElementById("maxSylCount") as HTMLInputElement;
const wordCountElement = document.getElementById("wordCount") as HTMLInputElement;
const wordOutputTextArea = document.getElementById("outputText") as HTMLInputElement;

const allowDuplicatesElement = document.getElementById("allowDuplicates") as HTMLInputElement;
const sortOutputElement = document.getElementById("sortOutput") as HTMLInputElement;
const debugOutputElement = document.getElementById("debugOutput") as HTMLInputElement;
const separateSyllablesElement = document.getElementById("separateSyllables") as HTMLInputElement;
const forceWordLimitElement = document.getElementById("forceWordLimit") as HTMLInputElement;

const duplicateAlertElement = document.getElementById("duplicateAlert") as HTMLElement;
const rejectedAlertElement = document.getElementById("rejectedAlert") as HTMLElement;

submit?.addEventListener("click", () => {
    wordOutputTextArea.value = "";
    duplicateAlertElement.innerHTML = "";
    duplicateAlertElement.hidden = true;
    rejectedAlertElement.innerHTML = "";
    rejectedAlertElement.hidden = true;

    const wordCount = Number.parseInt(wordCountElement.value, 10);
    let categories: CategoryListing = new Map<string, Category>();
    let tokens: Token[] = [];
    let syllable: Syllable | ParseError;
    const rejects: string[] = [];
    let letters: string[] = [];

    let minSylCount = Number.parseInt(minSylCountElement.value, 10);
    let maxSylCount = Number.parseInt(maxSylCountElement.value, 10);

    if (maxSylCount < minSylCount) {
        minSylCountElement.value = maxSylCount.toString();
        minSylCount = maxSylCount;
    } else if (minSylCount > maxSylCount) {
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
        } else if (line.match(/reject:/)) {
            rejects.push(
                ...line.replaceAll("reject:", "")
                    .trim()
                    .split("|")
                    .map((s) => s
                        .trim()
                        // extract rejections in brackets
                        .replace(/{(.+)}/, (_, g1) => g1))
                    .filter((s) => s.length > 0),
            );
        } else if (line.match(/letters:/)) {
            letters = line.replaceAll("letters:", "").split(" ").filter((e) => e.length > 0);
        }
    });
    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `letters: ${letters.join(",")}\n`;
        wordOutputTextArea.value += `rejections: ${rejects.join(",")}\n`;
    }

    const maybeCats: CategoryListing = new Map<string, Category>();
    try {
        Array.from(categories).forEach(((nameCat) => {
            const cat = fillCategory(nameCat[0], categories);
            cat.setWeights();
            maybeCats.set(nameCat[0], cat);
        }));
    } catch (e: any) {
        wordOutputTextArea.value = e;
        return;
    }
    categories = maybeCats;

    let rejectComps: Reject[] = [];
    let rejectRegexp: RegExp;
    if (rejects.length > 0) {
        try {
            rejectComps = rejects.map((r) => new Reject(r, categories));
            rejectRegexp = new RegExp(rejectComps.map((r) => r.toRegex().source).join("|"));
        } catch (e: any) {
            wordOutputTextArea.value = e;
            return;
        }
    } else {
        rejectRegexp = /$^/;
    }

    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `reject syls: ${rejectComps.map((r) => JSON.stringify(r)).join(" ")}\n`;
        wordOutputTextArea.value += `reject regex: ${rejectRegexp}\n`;
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
            return;
        }

        const possibleSyllableCount = syllable.evaluateAll().length;

        let words: string[][] = [];

        let rejectedCount = 0;
        let duplicateCount = 0;
        let generatedWords = 0;

        while (words.length < wordCount) {
            const syls = generateWord(syllable, minSylCount, maxSylCount);
            generatedWords += 1;

            if (rejectRegexp.test(syls.join(""))) {
                // rejections
                rejectedCount += 1;
            } else if (!allowDuplicatesElement.checked && [...new Set(words.map((s) => s.join("")))].includes(syls.join(""))) {
                // duplicates
                duplicateCount += 1;
            } else {
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
        rejectedAlertElement.innerHTML += `generated ${generatedWords} words, `;
        rejectedAlertElement.hidden = false;

        if (rejectedCount > 0) {
            rejectedAlertElement.innerHTML += `rejected ${rejectedCount} words`;
            rejectedAlertElement.hidden = false;
        }

        if (duplicateCount > 0) {
            duplicateAlertElement.innerHTML += `removed ${duplicateCount} duplicates`;
            duplicateAlertElement.hidden = false;
        }
        if (sortOutputElement.checked && letters.length > 0) {
            // sort based on letters
            // letters can be of any length
            // tokenize the words into their letters
            const letterizedWords = words.map(
                (w) => ({ word: w, lets: letterizeWord(w, letters) }),
            );
            // sort based on these letters
            words = letterizedWords.slice().sort(
                (left, right) => {
                    // convert to indexes per letter
                    const letterIndexer = toIndexArray(letters);
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
                },
            ).map((obj) => obj.word);
        } else if (sortOutputElement.checked) {
            words = words.slice().sort();
        }

        let outWords = words.map((syls) => syls.join(separateSyllablesElement.checked ? "." : ""));
        if (!allowDuplicatesElement.checked) {
            const wordset = [...new Set(outWords)];
            if (wordset.length < outWords.length) {
                duplicateAlertElement.innerHTML += `removed ${outWords.length - wordset.length} duplicates`;
                duplicateAlertElement.hidden = false;
            }
            outWords = wordset;
        }
        wordOutputTextArea.value += outWords.join("\n");
    }
});

// generate a word as its syllables
function generateWord(syllable: Syllable, minSyllables: number, maxSyllables: number): string[] {
    const outWord: string[] = [];
    const numSyllables = Math.max(
        minSyllables,
        Math.floor(maxSyllables - Math.random() * maxSyllables) + 1,
    );
    for (let i = 0; i < numSyllables; i += 1) {
        outWord.push(syllable.evaluate());
    }
    return outWord;
}

function toIndexArray(letters: string[]): (letters: string[]) => number[] {
    return (wordLetters: string[]): number[] => wordLetters.map((l) => letters.indexOf(l));
}

// tokenize a word (as its syllables) into a list of contained letters
function letterizeWord(word: string[], letters: string[]): string[] {
    return word.flatMap((syl) => letterizeSyllable(syl, letters));
}

// tokenize a syllable into letters
function letterizeSyllable(syllable: string, letters: string[]): string[] {
    const letterRegexp = new RegExp(`(${letters.slice().sort((a, b) => b.length - a.length).join("|")})`, "u");

    return syllable.split(letterRegexp).filter((s) => s.length > 0);
}

class Reject {
    matchWordStart: boolean;

    matchWordEnd: boolean;

    matchSylStart: boolean;

    matchSylEnd: boolean;

    rejectSyllable: Syllable;

    constructor(rejection: string, categories: CategoryListing) {
        this.matchWordStart = rejection.startsWith("^");
        this.matchWordEnd = rejection.endsWith(";");
        this.matchSylStart = rejection.startsWith("@");
        this.matchSylEnd = rejection.endsWith("&");
        const s = parseSyllable(tokenizeSyllable(rejection.replaceAll(/[&^;@]/g, "")), categories, rejection.replaceAll(/[&^;@]/g, ""));
        if (s instanceof ParseError) {
            throw s;
        }
        this.rejectSyllable = s;
    }

    toRegex(): RegExp {
        let reg = this.rejectSyllable.toRegex().source;
        if (this.matchWordStart) {
            reg = `^${reg}`;
        }
        if (this.matchWordEnd) {
            reg = `${reg}$`;
        }
        return new RegExp(`(${reg})`);
    }
}
