import CopiableTextArea from "./CopiableTextArea";

interface OutputProps {
    output: string;
    generatedCount: number;
    duplicateCount: number;
    rejectedCount: number;
    replacedCount: number;
}

const WAITING_FOR_INPUT = "Waiting for input...";

export default function Output(props: OutputProps) {
    return (
        <section className="col">
            <label htmlFor="outputText">
                <h2>Output</h2>
            </label>
            <div className="floating-container">
                <CopiableTextArea
                    id="outputText"
                    placeholder={props.output !== "" ? props.output : WAITING_FOR_INPUT}
                    isReadOnly={true}
                />
            </div>
            {props.generatedCount > -1 ? (
                <div className="alert alert-primary preline-whitespace" id="generatedAlert">
                    Generated {props.generatedCount} words.
                </div>
            ) : null}
            {props.duplicateCount > -1 ? (
                <div className="alert alert-info preline-whitespace" id="duplicateAlert">
                    Removed {props.duplicateCount} duplicates.
                </div>
            ) : null}
            {props.rejectedCount > -1 ? (
                <div className="alert alert-warning preline-whitespace" id="rejectedAlert">
                    Rejected {props.rejectedCount} words.
                </div>
            ) : null}
            {props.replacedCount > -1 ? (
                <div className="alert alert-secondary preline-whitespace" id="replacedAlert">
                    Replaced {props.replacedCount} words.
                </div>
            ) : null}
        </section>
    );
}
