export default function Changelog() {
    return (
        <section className="row">
            <h2>Changelog</h2>
            <ul>
                <li>
                    v1.4: "The React Update"
                    <ul>
                        <li>Decouple backend code from the frontend UI.</li>
                        <li>Migrate frontend UI to React.</li>
                    </ul>
                </li>
                <li>
                    v1.3.1:
                    <ul>
                        <li>
                            Maximum syllable count updates to be minimum syllable count if minimum is larger than
                            maximum.
                        </li>
                        <li>Fix "Force word limit" calculating an incorrect number of possibilites.</li>
                        <li>Now shows more errors if mutliple errors exist.</li>
                    </ul>
                </li>
                <li>v1.3: Add examples.</li>
                <li>v1.2: Allow any non-raw-ending symbol as a phoneme.</li>
                <li>v1.1.1: Fix named components having choice count of 0.</li>
                <li>v1.1.0: Add named components.</li>
                <li>v1.0.0: Go rewrite release. Word, sentence gen; rejections; sorting; syllable marking, etc.</li>
            </ul>
        </section>
    );
}
