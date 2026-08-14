interface OutputProps {
    output: string;
    generatedCount: number;
    duplicateCount: number;
    rejectedCount: number;
    replacedCount: number;
}

export default function Output(props: OutputProps) {
    return (
        <section className="col">
            <label htmlFor="outputText">
                <h3>Output</h3>
            </label>
            <textarea
                id="outputText"
                className="form-control"
                rows={25}
                readOnly
                defaultValue="Waiting for input..."
            ></textarea>
            <button id="copyButton" className="btn btn-danger">
                Copy
            </button>
            {props.generatedCount ? (
                <div className="alert alert-primary preline-whitespace" id="generatedAlert">
                    Waiting for input...
                </div>
            ) : null}
            {props.duplicateCount ? (
                <div className="alert alert-info preline-whitespace" id="duplicateAlert">
                    Waiting for input...
                </div>
            ) : null}
            {props.rejectedCount ? (
                <div className="alert alert-warning preline-whitespace" id="rejectedAlert">
                    Waiting for input...
                </div>
            ) : null}
            {props.replacedCount ? (
                <div className="alert alert-secondary preline-whitespace" id="replacedAlert">
                    Waiting for input...
                </div>
            ) : null}
        </section>
    );
}
