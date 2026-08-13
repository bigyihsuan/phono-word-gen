export default function Output() {
    return (
        <section>
            <label htmlFor="outputText">
                <h3>Output</h3>
            </label>
            <textarea id="outputText" className="form-control" rows={25} readOnly></textarea>
            <button id="copyButton" className="btn btn-danger">
                Copy
            </button>
            <div className="alert alert-primary preline-whitespace" id="generatedAlert">
                Waiting for input...
            </div>
            <div className="alert alert-info preline-whitespace" id="duplicateAlert">
                Waiting for input...
            </div>
            <div className="alert alert-warning preline-whitespace" id="rejectedAlert">
                Waiting for input...
            </div>
            <div className="alert alert-secondary preline-whitespace" id="replacedAlert">
                Waiting for input...
            </div>
        </section>
    );
}
