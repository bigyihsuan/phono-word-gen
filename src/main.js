var phonology = document.getElementById("phonology");
var syllable = document.getElementById("syllable");
var submit = document.getElementById("submit");
submit === null || submit === void 0 ? void 0 : submit.addEventListener("click", function () {
    console.log(phonology === null || phonology === void 0 ? void 0 : phonology.value);
    console.log(syllable === null || syllable === void 0 ? void 0 : syllable.value);
});
