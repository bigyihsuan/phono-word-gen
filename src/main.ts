const phonology = document.getElementById("phonology") as HTMLInputElement
const syllable = document.getElementById("syllable") as HTMLInputElement
const submit = document.getElementById("submit") as HTMLButtonElement

submit?.addEventListener("click", function () {
    console.log(phonology?.value)
    console.log(syllable?.value)
})