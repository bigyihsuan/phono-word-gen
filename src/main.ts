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

submit?.addEventListener("click", () => {
    wordOutputTextArea.value = "";

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
                    .split(/[{}]/)
                    .map((s) => s
                        .trim()
                        .replace(/{(.+)}/, (_, g1) => g1))
                    .filter((s) => s.length > 0),
            );
        }
    });
    if (debugOutputElement.checked) {
        wordOutputTextArea.value += `${rejects.join(",")}\n`;
    }

    let maybeCats: CategoryListing;
    try {
        maybeCats = fillCategories(categories);
    } catch (e: any) {
        wordOutputTextArea.value = e;
        return;
    }
    categories = maybeCats;

    const sylLine = lines.find((l) => l.trim().match(/syllable:/))?.replaceAll("syllable:", "").trim();
    if (sylLine !== undefined) {
        tokens = tokenizeSyllable(sylLine);
        syllable = parseSyllable(tokens.slice(), categories, sylLine);
        if (debugOutputElement.checked) {
            wordOutputTextArea.value += syllable.toString();
            wordOutputTextArea.value += "\n---------------\n";
        }

        if (syllable instanceof ParseError) {
            wordOutputTextArea.value += syllable.toString();
        } else if (syllable !== undefined) {
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
                words = [...new Set(words)];
            }
            if (sortOutputElement.checked) {
                words = words.sort();
            }
            wordOutputTextArea.value += words.join("\n");
        }
    }
});
