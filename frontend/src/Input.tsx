import { useState } from "react";
import CopiableTextArea from "./CopiableTextArea";

enum GenerateType {
    words = "words",
    sentences = "sentences",
}

interface InputState {
    phonology: string;
    minSylCount: number;
    maxSylCount: number;
    wordCount: number;
    sentenceCount: number;
    generateType: GenerateType;
    forbidDuplicates: boolean;
    forceWordLimit: boolean;
    applyRejections: boolean;
    applyReplacements: boolean;
    markSyllables: boolean;
    sortOutput: boolean;
    debugOutput: boolean;
}

export default function Input() {
    const placeholder = `# This is a simple example phonology. You can type in this box!
C = p t k
V = a i u
syllable: $C$V($C)*25`;

    const [input, setInput] = useState({
        phonology: "",
        minSylCount: 1,
        maxSylCount: 5,
        wordCount: 25,
        sentenceCount: 5,
        generateType: GenerateType.words,
        forbidDuplicates: false,
        forceWordLimit: false,
        applyRejections: false,
        applyReplacements: false,
        markSyllables: false,
        sortOutput: false,
        debugOutput: false,
    } as InputState);

    // This is unfortunately a controlled component
    // because <form action> empties out the form.
    // Then I would need to fill in the values into the form again,
    // but then at that point I might as well go all the way and make it a controlled component.

    //#region onChange handlers
    function phonologyChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            phonology: (e.target as HTMLTextAreaElement).value,
        });
    }

    function minSylCountChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            minSylCount: Number.parseInt((e.target as HTMLInputElement).value),
        });
    }
    function maxSylCountChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            maxSylCount: Number.parseInt((e.target as HTMLInputElement).value),
        });
    }
    function wordCountChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            wordCount: Number.parseInt((e.target as HTMLInputElement).value),
        });
    }
    function sentenceCountChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            sentenceCount: Number.parseInt((e.target as HTMLInputElement).value),
        });
    }

    function generateTypeChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            generateType: (e.target as HTMLInputElement).value as GenerateType,
        });
    }
    function forbidDuplicatesChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            forbidDuplicates: (e.target as HTMLInputElement).value === "on",
        });
    }
    function forceWordLimitChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            forceWordLimit: (e.target as HTMLInputElement).value === "on",
        });
    }
    function applyRejectionsChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            applyRejections: (e.target as HTMLInputElement).value === "on",
        });
    }
    function applyReplacementsChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            applyReplacements: (e.target as HTMLInputElement).value === "on",
        });
    }
    function markSyllablesChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            markSyllables: (e.target as HTMLInputElement).value === "on",
        });
    }
    function sortOutputChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            sortOutput: (e.target as HTMLInputElement).value === "on",
        });
    }
    function debugOutputChanges(e: React.ChangeEvent) {
        setInput({
            ...input,
            debugOutput: (e.target as HTMLInputElement).value === "on",
        });
    }
    //#endregion onChange handlers

    function handleSubmit(event: React.MouseEvent) {
        event.preventDefault();
        console.log(input);
        // TODO: hook WASM into this
        // TODO: give output to output component
    }

    return (
        <section className="col">
            <label htmlFor="inputs">
                <h2>Input</h2>
            </label>
            <form id="inputs" className="input">
                <CopiableTextArea
                    id="phonology"
                    placeholder={placeholder}
                    isReadOnly={false}
                    onChange={phonologyChanges}
                />
                <div className="form-control">
                    <div className="input-group">
                        <span className="input-group-text">Min Syllables:</span>
                        <input
                            type="number"
                            min={1}
                            value={input.minSylCount}
                            id="minSylCount"
                            name="minSylCount"
                            onChange={minSylCountChanges}
                            className="form-control"
                        />
                        <span className="input-group-text">Max Syllables:</span>
                        <input
                            type="number"
                            min={1}
                            value={input.maxSylCount}
                            id="maxSylCount"
                            name="maxSylCount"
                            onChange={maxSylCountChanges}
                            className="form-control"
                        />
                    </div>
                    {input.generateType === GenerateType.words ? (
                        <div className="input-group" id="wordCountInput">
                            <span className="input-group-text">Number of Words</span>
                            <input
                                type="number"
                                min={1}
                                value={input.wordCount}
                                id="wordCount"
                                name="wordCount"
                                onChange={wordCountChanges}
                                className="form-control"
                            />
                        </div>
                    ) : (
                        <div className="input-group" id="sentenceCountInput">
                            <span className="input-group-text">Number of Sentences</span>
                            <input
                                type="number"
                                min={1}
                                value={input.sentenceCount}
                                id="sentenceCount"
                                name="sentenceCount"
                                onChange={sentenceCountChanges}
                                className="form-control"
                            />
                        </div>
                    )}
                </div>
                <div className="form-control container">
                    <div className="container">
                        <div className="input-group row">
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="radio"
                                    name="generateType"
                                    id="generateWords"
                                    value={GenerateType.words}
                                    onChange={generateTypeChanges}
                                    defaultChecked
                                />
                                <label className="form-check-label" htmlFor="generateWords">
                                    Generate words
                                </label>
                            </div>
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="radio"
                                    name="generateType"
                                    id="generateSentences"
                                    value={GenerateType.sentences}
                                    onChange={generateTypeChanges}
                                />
                                <label className="form-check-label" htmlFor="generateSentences">
                                    Generate sentences
                                </label>
                            </div>
                        </div>
                        <div className="input-group row">
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="checkbox"
                                    id="forbidDuplicates"
                                    name="forbidDuplicates"
                                    onChange={forbidDuplicatesChanges}
                                />
                                <label className="form-check-label" htmlFor="forbidDuplicates">
                                    Forbid duplicates
                                </label>
                            </div>
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="checkbox"
                                    id="forceWordLimit"
                                    name="forceWordLimit"
                                    onChange={forceWordLimitChanges}
                                />
                                <label className="form-check-label" htmlFor="forceWordLimit">
                                    Force word limit
                                </label>
                            </div>
                        </div>
                        <div className="input-group row">
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="checkbox"
                                    id="applyRejections"
                                    name="applyRejections"
                                    onChange={applyRejectionsChanges}
                                />
                                <label className="form-check-label" htmlFor="applyRejections">
                                    Apply rejections
                                </label>
                            </div>
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="checkbox"
                                    id="applyReplacements"
                                    name="applyReplacements"
                                    onChange={applyReplacementsChanges}
                                    disabled
                                />
                                <label className="form-check-label" htmlFor="applyReplacements">
                                    Apply replacements
                                </label>
                            </div>
                        </div>
                        <div className="input-group row">
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="checkbox"
                                    id="markSyllables"
                                    name="markSyllables"
                                    onChange={markSyllablesChanges}
                                />
                                <label className="form-check-label" htmlFor="markSyllables">
                                    Mark syllables
                                </label>
                            </div>
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="checkbox"
                                    id="sortOutput"
                                    name="sortOutput"
                                    onChange={sortOutputChanges}
                                />
                                <label className="form-check-label" htmlFor="sortOutput">
                                    Sort output
                                </label>
                            </div>
                            <div className="form-check form-check-inline">
                                <input
                                    className="form-check-input"
                                    type="checkbox"
                                    id="debugOutput"
                                    name="debugOutput"
                                    onChange={debugOutputChanges}
                                />
                                <label className="form-check-label" htmlFor="debugOutput">
                                    Include debug output
                                </label>
                            </div>
                        </div>
                    </div>
                    <button className="input-group btn btn-primary" onClick={handleSubmit}>
                        Generate {input.generateType}
                    </button>
                </div>
            </form>
        </section>
    );
}
