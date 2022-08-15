import {
    Category, CategoryListing, fillCategories, parseCategory,
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
                    // extract rejections in curly brackets
                    .split(/[<>]/)
                    .map((s) => s
                        .trim()
                        .replace(/{(.+)}/, (_, g1) => g1))
                    .filter((s) => s.length > 0),
            );
        }
    });
    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `rejections: ${rejects.join(",")}\n`;
    }

    let maybeCats: CategoryListing;
    try {
        maybeCats = fillCategories(categories);
        maybeCats.forEach((cat) => cat.setWeights());
    } catch (e: any) {
        wordOutputTextArea.value = e;
        return;
    }
    categories = maybeCats;

    const rejectComps = rejects.map(
        (r) => parseSyllable(tokenizeSyllable(r), categories, r),
    ).filter((r) => r instanceof Syllable);

    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `parsed rejects:\n${rejectComps.join("\n")}\n`;
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
        } else {
            if (debugOutputElement.checked) {
                wordOutputTextArea.value += `possibles: ${syllable.possibilities}`;
                wordOutputTextArea.value += "\n---------------\n";
            }
            let words: string[] = [];
            for (let _ = 0; _ < wordCount; _ += 1) {
                let outWord = "";
                const numSyllables = Math.max(
                    minSylCount,
                    Math.floor(maxSylCount - Math.random() * maxSylCount) + 1,
                );
                for (let i = 0; i < numSyllables; i += 1) {
                    outWord += syllable.evaluate();
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
            const withoutRejected = words.filter(
                (w: string) => !rejectComps.some(
                    (r) => (r instanceof Syllable) && r.matches(w),
                ),
            );
            if (words.length > withoutRejected.length) {
                rejectedAlertElement.innerHTML += `rejected ${words.length - withoutRejected.length} words`;
                rejectedAlertElement.hidden = false;
                words = withoutRejected;
            }

            if (sortOutputElement.checked) {
                words = words.sort();
            }
            wordOutputTextArea.value += words.join("\n");
        }
    }
});
