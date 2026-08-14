import { useEffect, useState } from "react";

interface CopiableTextAreaProps {
    id: string;
    placeholder: string;
    isReadOnly: boolean;
}

export default function CopiableTextArea(props: CopiableTextAreaProps) {
    const [value, setValue] = useState("");

    useEffect(() => {
        const target = document.getElementById(props.id)! as HTMLTextAreaElement;
        function watchTarget() {
            setValue(target.value);
        }
        target.addEventListener("change", watchTarget);
        return () => {
            target.removeEventListener("change", watchTarget);
        };
    }, []);

    function makeCopy() {
        window.navigator.clipboard.writeText(value);
    }

    return (
        <>
            <textarea
                id={props.id}
                className="form-control"
                rows={25}
                placeholder={props.placeholder}
                readOnly={props.isReadOnly}
            ></textarea>
            <button type="button" className="btn btn-secondary floating-bottom-right" onClick={makeCopy}>
                Copy
            </button>
        </>
    );
}
