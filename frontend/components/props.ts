export interface InputProps {
    inputState: InputState;
    handleSubmit: (event: React.MouseEvent) => void;
    phonologyChanges: (event: React.ChangeEvent) => void;
    minSylCountChanges: (event: React.ChangeEvent) => void;
    maxSylCountChanges: (event: React.ChangeEvent) => void;
    wordCountChanges: (event: React.ChangeEvent) => void;
    sentenceCountChanges: (event: React.ChangeEvent) => void;
    generateTypeChanges: (event: React.ChangeEvent) => void;
    forbidDuplicatesChanges: (event: React.ChangeEvent) => void;
    forceWordLimitChanges: (event: React.ChangeEvent) => void;
    applyRejectionsChanges: (event: React.ChangeEvent) => void;
    applyReplacementsChanges: (event: React.ChangeEvent) => void;
    markSyllablesChanges: (event: React.ChangeEvent) => void;
    sortOutputChanges: (event: React.ChangeEvent) => void;
    debugOutputChanges: (event: React.ChangeEvent) => void;
}

export interface InputState {
    phonology: string;
    minSylCount: number;
    maxSylCount: number;
    wordCount: number;
    sentenceCount: number;
    generateType: "words" | "sentences";
    forbidDuplicates: boolean;
    forceWordLimit: boolean;
    applyRejections: boolean;
    applyReplacements: boolean;
    markSyllables: boolean;
    sortOutput: boolean;
    debugOutput: boolean;
}

export interface OutputState {
    output: string;
    generatedCount: number;
    duplicateCount: number;
    rejectedCount: number;
    replacedCount: number;
}
