import { parseCategory } from "./modules/category.js";
import { tokenizeSyllable } from "./modules/syllable.js";
const phonology = document.getElementById("phonology");
const submit = document.getElementById("submit");
const wordOutput = document.getElementById("output");
let categories = new Map();
let tokens;
submit === null || submit === void 0 ? void 0 : submit.addEventListener("click", function () {
    let lines = phonology === null || phonology === void 0 ? void 0 : phonology.value.replaceAll(/\n+/g, "\n").split("\n").filter((s) => s.length > 0);
    lines.forEach((l) => {
        let line = l.trim();
        if (line.match(/=/)) {
            let cat = parseCategory(line);
            categories.set(cat.name, cat);
        }
        else if (line.match(/syllable:/)) {
            // TODO: parse syllable structure
            tokens = tokenizeSyllable(line);
        }
    });
    wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.replaceChildren(); // clear the output for each run
    categories.forEach((c) => {
        let p = document.createElement("p");
        p.innerHTML = c.toString();
        wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.append(p);
    });
    wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.append(document.createElement("hr"));
    tokens.forEach((t) => {
        let p = document.createElement("p");
        p.innerHTML = t.toString();
        wordOutput === null || wordOutput === void 0 ? void 0 : wordOutput.append(p);
    });
});
//# sourceMappingURL=main.js.map