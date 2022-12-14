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
    const replStrs: string[] = [];

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

    const letterwiseCompare = compareWordsLetterwise(letters);

    const maybeCats: CategoryListing = new Map<string, Category>();
    try {
        Array.from(categories).forEach(((nameCat) => {
            const cat = fillCategory(nameCat[0], categories);
            cat.setWeights();
            maybeCats.set(nameCat[0], cat);
        }));
    } catch (e: any) {
        wordOutputTextArea.value = e;
        console.error(e);
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
            console.error(e);
            return;
        }
    } else {
        rejectRegexp = /$^/;
    }

    const replacements: Replacement[] = [];
    try {
        replStrs.forEach((r) => {
            replacements.push(new Replacement(r, categories));
        });
    } catch (e: any) {
        wordOutputTextArea.value += e;
        console.error(e);
        return;
    }
    // replacements.forEach((r) => {
    //     console.log(r);
    // });

    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `letters: ${letters.join(",")}\n\n`;
        wordOutputTextArea.value += `rejections: ${rejects.join(",")}\n\n`;
        wordOutputTextArea.value += `replacements:\n    ${replacements.map((r) => r.toString()).join("\n    ")}\n\n`;
        wordOutputTextArea.value += `reject regex: ${rejectRegexp}\n\n`;
        wordOutputTextArea.value += `categories: ${Array.from(categories).map((cn) => cn[1]).join("\n")}\n\n`;
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
            console.error(syllable);
            return;
        }

        const possibleSyllableCount = syllable.evaluateAll().length;

        let words: string[][] = [];

        let rejectedCount = 0;
        let duplicateCount = 0;
        let generatedWords = 0;
        let replacedWords = 0;

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
        if (sortOutputElement.checked && letters.length > 0) {
            // sort based on letters
            // letters can be of any length
            // tokenize the words into their letters
            const letterizedWords: letterizedWord[] = words.map(
                (w) => ({ word: w, lets: letterizeWord(w, letters) }),
            );
            // sort based on these letters
            const compare = compareWordsLetterwise(letters);
            words = letterizedWords.slice().sort(compare).map((obj) => obj.word);
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
                replacedWords += 1;
            }
            return w;
        });

        // resort after doing replacements
        if (sortOutputElement.checked) {
            const letterizedWords: letterizedWord[] = outWords.map(
                (w) => ({ word: [w], lets: letterizeWord([w], letters) }),
            );
            outWords = letterizedWords.slice().sort(letterwiseCompare).map((obj) => obj.word).map((sarr) => sarr.join(""));
        }

        rejectedAlertElement.innerHTML += `, replaced ${replacedWords} words`;
        rejectedAlertElement.hidden = false;

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

type letterizedWord = { word: string[], lets: string[] };

function compareWordsLetterwise(letters: string[]):
    (left: letterizedWord, right: letterizedWord) => number {
    // convert to indexes per letter

    return (left: letterizedWord, right: letterizedWord) => {
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
    };
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
