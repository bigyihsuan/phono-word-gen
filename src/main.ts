import { Category, parseCategory } from "./modules/category.js"
import { tokenizeSyllable, Token } from "./modules/syllable.js"

const phonology = document.getElementById("phonology") as HTMLInputElement
const submit = document.getElementById("submit") as HTMLButtonElement
const wordOutput = document.getElementById("output")

let categories: Map<string, Category> = new Map<string, Category>();
let tokens: Token[]

submit?.addEventListener("click", function () {
    let lines = phonology?.value.replaceAll(/\n+/g, "\n").split("\n").filter((s) => s.length > 0);

    lines.forEach((l) => {
        let line = l.trim()
        if (line.match(/=/)) {
            let cat = parseCategory(line);
            categories.set(cat.name, cat);
        } else if (line.match(/syllable:/)) {
            // TODO: parse syllable structure
            tokens = tokenizeSyllable(line)
        }
    })

    wordOutput?.replaceChildren(); // clear the output for each run

    categories.forEach((c) => {
        let p: HTMLParagraphElement = document.createElement("p")
        p.innerHTML = c.toString()
        wordOutput?.append(p)
    })
    wordOutput?.append(document.createElement("hr"))
    tokens.forEach((t) => {
        let p: HTMLParagraphElement = document.createElement("p")
        p.innerHTML = t.toString()
        wordOutput?.append(p)
    })
})

