import { useEffect, useState } from "react";

interface CopyButtonProps {
    textAreaId: string;
}

export default function CopyButton(props: CopyButtonProps) {
    const [value, setValue] = useState("");

    useEffect(() => {
        const target = document.getElementById(props.textAreaId)! as HTMLTextAreaElement;
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
        <button type="button" className="btn btn-secondary floating-bottom-right" onClick={makeCopy}>
            Copy
        </button>
    );
}
