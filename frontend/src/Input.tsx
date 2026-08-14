export default function Input() {
    const defaultInput = `# This is a simple example phonology. You can type in this box!
C = p t k
V = a i u
syllable: $C$V($C)*25`;

    return (
        <section className="col">
            <label htmlFor="inputs">
                <h2>Input</h2>
            </label>
            <form id="inputs" className="input">
                <div className="floating-container">
                    <textarea id="phonology" className="form-control" rows={25} placeholder={defaultInput}></textarea>
                    <button id="copyInputButton" className="btn btn-secondary floating-bottom-right">
                        Copy
                    </button>
                </div>
                <div id="numberInputs" className="form-control">
                    <div className="input-group" id="syllableCountInput">
                        <span className="input-group-text">Min/Max syllables/word</span>
                        <input type="number" min="1" defaultValue="1" id="minSylCount" className="form-control" />
                        <span className="input-group-text">&ndash;</span>
                        <input type="number" min="1" defaultValue="1" id="maxSylCount" className="form-control" />
                    </div>
                    <div className="input-group" id="wordCountInput">
                        <span className="input-group-text">Number of Words</span>
                        <input type="number" min="1" defaultValue="25" id="wordCount" className="form-control" />
                    </div>
                    <div className="input-group" id="sentenceCountInput">
                        <span className="input-group-text">Number of Sentences</span>
                        <input
                            type="number"
                            min="1"
                            defaultValue="5"
                            id="sentenceCount"
                            className="form-control"
                            disabled
                        />
                    </div>
                </div>
                <div id="checkboxes" className="form-control container row">
                    <div className="container">
                        <div className="input-group row">
                            <div className="form-check form-check-inline col">
                                <input
                                    className="form-check-input"
                                    type="radio"
                                    name="generateType"
                                    id="generateWords"
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
                                />
                                <label className="form-check-label" htmlFor="generateSentences">
                                    Generate sentences
                                </label>
                            </div>
                        </div>
                        <div className="input-group row">
                            <div className="form-check form-check-inline col">
                                <input className="form-check-input" type="checkbox" id="forbidDuplicates" />
                                <label className="form-check-label" htmlFor="forbidDuplicates">
                                    Forbid duplicates
                                </label>
                            </div>
                            <div className="form-check form-check-inline col">
                                <input className="form-check-input" type="checkbox" id="forceWordLimit" />
                                <label className="form-check-label" htmlFor="forceWordLimit">
                                    Force word limit
                                </label>
                            </div>
                        </div>
                        <div className="input-group row">
                            <div className="form-check form-check-inline col">
                                <input className="form-check-input" type="checkbox" id="applyRejections" />
                                <label className="form-check-label" htmlFor="applyRejections">
                                    Apply rejections
                                </label>
                            </div>
                            <div className="form-check form-check-inline col">
                                <input className="form-check-input" type="checkbox" id="applyReplacements" disabled />
                                <label className="form-check-label" htmlFor="applyReplacements">
                                    Apply replacements
                                </label>
                            </div>
                        </div>
                        <div className="input-group row">
                            <div className="form-check form-check-inline col">
                                <input className="form-check-input" type="checkbox" id="markSyllables" />
                                <label className="form-check-label" htmlFor="markSyllables">
                                    Mark syllables
                                </label>
                            </div>
                            <div className="form-check form-check-inline col">
                                <input className="form-check-input" type="checkbox" id="sortOutput" />
                                <label className="form-check-label" htmlFor="sortOutput">
                                    Sort output
                                </label>
                            </div>
                            {/* <div className="form-check form-check-inline">
                                <input className="form-check-input" type="checkbox" id="debugOutput" />
                                <label className="form-check-label" htmlFor="debugOutput">
                                    Include debug output
                                </label>
                            </div> */}
                        </div>
                    </div>
                    <button type="button" id="submit" className="btn btn-primary col">
                        Generate
                    </button>
                </div>
            </form>
        </section>
    );
}
