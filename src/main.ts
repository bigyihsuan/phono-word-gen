const phonology = document.getElementById("phonology") as HTMLInputElement
const submit = document.getElementById("submit") as HTMLButtonElement

submit?.addEventListener("click", function () {
    let lines = phonology?.value.split("\n");
    for (let idx in lines) {
        let line = lines[idx].trim()
        console.log(line)
        if (line.match(/=/)) {
            console.log(parseCategory(line))
        }
    }
})

type Phoneme = string;

class Category {
    name: string;
    phonemes: Phoneme[];

    constructor(name: string, phonemes: Phoneme[]) {
        this.name = name;
        this.phonemes = phonemes;
    }
}

function parseCategory(cat: string): Category {
    let name: string = "";
    let phonemes: Phoneme[] = [];
    let split = cat.trim().split("=").map((s) => s.trim()); // split on the equals and trim both sides

    name = split[0];
    phonemes = split[1].split(" ")

    return new Category(name, phonemes);
}
