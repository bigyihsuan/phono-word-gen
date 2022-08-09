const stuff = document.getElementById("stuff")

const p = document.createElement("p")

p.textContent = "hello world!"

stuff?.appendChild(p)