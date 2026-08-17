import CopiableTextArea from "./CopiableTextArea";
import type { InputProps } from "./props";

export default function Input(props: InputProps) {
    const placeholder = `# This is a simple example phonology. You can type in this box!
C = p t k
V = a i u
syllable: $C$V($C)*25`;

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
                    onChange={props.phonologyChanges}
                />
                <div className="form-control">
                    <div className="input-group">
                        <span className="input-group-text">Min Syllables:</span>
                        <input
                            type="number"
                            min={1}
                            value={props.inputState.minSylCount}
                            id="minSylCount"
                            name="minSylCount"
                            onChange={props.minSylCountChanges}
                            className="form-control"
                        />
                        <span className="input-group-text">Max Syllables:</span>
                        <input
                            type="number"
                            min={1}
                            value={props.inputState.maxSylCount}
                            id="maxSylCount"
                            name="maxSylCount"
                            onChange={props.maxSylCountChanges}
                            className="form-control"
                        />
                    </div>
                    {props.inputState.generateType === "words" ? (
                        <div className="input-group" id="wordCountInput">
                            <span className="input-group-text">Number of Words</span>
                            <input
                                type="number"
                                min={1}
                                value={props.inputState.wordCount}
                                id="wordCount"
                                name="wordCount"
                                onChange={props.wordCountChanges}
                                className="form-control"
                            />
                        </div>
                    ) : (
                        <div className="input-group" id="sentenceCountInput">
                            <span className="input-group-text">Number of Sentences</span>
                            <input
                                type="number"
                                min={1}
                                value={props.inputState.sentenceCount}
                                id="sentenceCount"
                                name="sentenceCount"
                                onChange={props.sentenceCountChanges}
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
                                    value="words"
                                    onChange={props.generateTypeChanges}
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
                                    value="sentences"
                                    onChange={props.generateTypeChanges}
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
                                    onChange={props.forbidDuplicatesChanges}
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
                                    onChange={props.forceWordLimitChanges}
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
                                    onChange={props.applyRejectionsChanges}
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
                                    onChange={props.applyReplacementsChanges}
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
                                    onChange={props.markSyllablesChanges}
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
                                    onChange={props.sortOutputChanges}
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
                                    onChange={props.debugOutputChanges}
                                />
                                <label className="form-check-label" htmlFor="debugOutput">
                                    Include debug output
                                </label>
                            </div>
                        </div>
                    </div>
                    <button className="input-group btn btn-primary" onClick={props.handleSubmit}>
                        Generate {props.inputState.generateType}
                    </button>
                </div>
            </form>
        </section>
    );
}
