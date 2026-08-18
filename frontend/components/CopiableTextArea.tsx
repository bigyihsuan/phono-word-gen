interface CopiableTextAreaProps {
    id: string;
    placeholder: string;
    value?: string;
    isReadOnly: boolean;
    onChange?: React.ChangeEventHandler;
    className?: string;
}

export default function CopiableTextArea(props: CopiableTextAreaProps) {
    function makeCopy() {
        const textArea = document.getElementById(props.id)! as HTMLTextAreaElement;
        window.navigator.clipboard.writeText(textArea.value);
    }

    return (
        <div className="floating-container">
            <textarea
                id={props.id}
                name={props.id}
                onChange={props.onChange}
                className={props.className !== undefined ? `form-control ${props.className}` : "form-control"}
                rows={25}
                placeholder={props.placeholder}
                value={props.value}
                readOnly={props.isReadOnly}
            ></textarea>
            <button type="button" className="btn btn-secondary floating-bottom-right" onClick={makeCopy}>
                Copy
            </button>
        </div>
    );
}
