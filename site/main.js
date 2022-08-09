"use strict";
const phonology = document.getElementById("phonology");
const submit = document.getElementById("submit");
submit === null || submit === void 0 ? void 0 : submit.addEventListener("click", function () {
    console.log(parseCategory(phonology === null || phonology === void 0 ? void 0 : phonology.value));
    let lines = phonology === null || phonology === void 0 ? void 0 : phonology.value.split("\n");
    for (let idx in lines) {
        let line = lines[idx].trim();
        console.log(line);
        if (line.match(/=/)) {
            console.log(parseCategory(line));
        }
    }
});
class Category {
    constructor(name, phonemes) {
        this.name = name;
        this.phonemes = phonemes;
    }
}
function parseCategory(cat) {
    let name = "";
    let phonemes = [];
    let split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides
    name = split[0];
    phonemes = split[1].split(" ");
    return new Category(name, phonemes);
}
//# sourceMappingURL=main.js.map