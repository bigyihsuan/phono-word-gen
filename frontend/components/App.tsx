import { useState } from "react";
import Footer from "./Footer";
import Header from "./Header";
import Input from "./Input";
import Output from "./Output";
import QuickReference from "./QuickDocs";
import type { InputState, OutputData } from "./props";

declare namespace globalThis {
    function generate(input: InputState): GenerateOutput;
}

interface GenerateOutput {
    words: string;
    sentences: string;
    errors: string;
    generatedCount: number;
    duplicateCount: number;
    rejectedCount: number;
    replacedCount: number;
}

export default function App() {
    const [inputState, setInputState] = useState<InputState>({
        phonology: "",
        minSylCount: 1,
        maxSylCount: 5,
        wordCount: 25,
        sentenceCount: 5,
        generateType: "words",
        forbidDuplicates: false,
        forceWordLimit: false,
        applyRejections: false,
        applyReplacements: false,
        markSyllables: false,
        sortOutput: false,
    });

    const outputStart = {
        output: "",
        generatedCount: -1,
        duplicateCount: -1,
        rejectedCount: -1,
        replacedCount: -1,
        hasError: false,
    };

    const [outputState, setOutputState] = useState<OutputData>({ ...outputStart });

    function handleSubmit(event: React.MouseEvent) {
        event.preventDefault();
        // call WASM function
        const obj = globalThis.generate(inputState);
        // give output to output component
        setOutputState({
            generatedCount: obj.generatedCount,
            duplicateCount: obj.duplicateCount,
            rejectedCount: obj.rejectedCount,
            replacedCount: obj.replacedCount,
            hasError: obj.errors !== "",
            output: obj.errors !== "" ? obj.errors : inputState.generateType === "words" ? obj.words : obj.sentences,
        });
    }

    /*
     * This is unfortunately a controlled component
     * because <form action> empties out the form.
     * Then I would need to fill in the values into the form again,
     * but then at that point I might as well go all the way and make it a controlled component.
     */

    //#region onChange handlers
    function phonologyChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, phonology: (e.target as HTMLTextAreaElement).value });
    }

    function minSylCountChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, minSylCount: Number.parseInt((e.target as HTMLInputElement).value) });
    }
    function maxSylCountChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, maxSylCount: Number.parseInt((e.target as HTMLInputElement).value) });
    }
    function wordCountChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, wordCount: Number.parseInt((e.target as HTMLInputElement).value) });
    }
    function sentenceCountChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, sentenceCount: Number.parseInt((e.target as HTMLInputElement).value) });
    }

    function generateTypeChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, generateType: (e.target as HTMLInputElement).value as "words" | "sentences" });
    }
    function forbidDuplicatesChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, forbidDuplicates: (e.target as HTMLInputElement).value === "on" });
    }
    function forceWordLimitChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, forceWordLimit: (e.target as HTMLInputElement).value === "on" });
    }
    function applyRejectionsChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, applyRejections: (e.target as HTMLInputElement).value === "on" });
    }
    function applyReplacementsChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, applyReplacements: (e.target as HTMLInputElement).value === "on" });
    }
    function markSyllablesChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, markSyllables: (e.target as HTMLInputElement).value === "on" });
    }
    function sortOutputChanges(e: React.ChangeEvent) {
        setInputState({ ...inputState, sortOutput: (e.target as HTMLInputElement).value === "on" });
    }
    //#endregion onChange handlers

    function loadExample(event: React.MouseEvent) {
        event.preventDefault();
        const exampleSelect = document.querySelector("#exampleSelect")! as HTMLSelectElement;
        const exampleName = exampleSelect.value;
        if (exampleName === "") {
            return;
        }
        console.log(exampleName);
        const fileName = `sample/${exampleName}.txt`;
        fetch(fileName)
            .then((response) => {
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                return response.text();
            })
            .then((exampleStr) => {
                setInputState((prev) => ({ ...prev, phonology: exampleStr })); // required trailing space
                clearOutput();
            });
    }

    function clearOutput() {
        setOutputState({ ...outputStart });
    }

    return (
        <>
            <Header />
            <main className="container">
                <div className="row">
                    <Input
                        inputState={inputState}
                        handleSubmit={handleSubmit}
                        phonologyChanges={phonologyChanges}
                        minSylCountChanges={minSylCountChanges}
                        maxSylCountChanges={maxSylCountChanges}
                        wordCountChanges={wordCountChanges}
                        sentenceCountChanges={sentenceCountChanges}
                        generateTypeChanges={generateTypeChanges}
                        forbidDuplicatesChanges={forbidDuplicatesChanges}
                        forceWordLimitChanges={forceWordLimitChanges}
                        applyRejectionsChanges={applyRejectionsChanges}
                        applyReplacementsChanges={applyReplacementsChanges}
                        markSyllablesChanges={markSyllablesChanges}
                        sortOutputChanges={sortOutputChanges}
                        loadExample={loadExample}
                    />
                    <Output
                        output={outputState.output}
                        generatedCount={outputState.generatedCount}
                        duplicateCount={outputState.duplicateCount}
                        rejectedCount={outputState.rejectedCount}
                        replacedCount={outputState.replacedCount}
                        hasError={outputState.hasError}
                    />
                </div>
                <hr></hr>
                <QuickReference />
            </main>
            <Footer />
        </>
    );
}
