export default function QuickDocs() {
    return (
        <section className="row">
            <h2>Quick Docs</h2>
            <h3>General</h3>
            <ul>
                <li>
                    Comments: <code>{`# comment ends at the end of the line`}</code>
                </li>
            </ul>
            <h3>Weighting Rules</h3>
            <ul>
                <li>Weights mark how often a component or phoneme can appear.</li>
                <li>Weights are positive integers.</li>
                <li>Weights can be applied to phonemes, optionals, or elements in a selection.</li>
                <li>
                    Phoneme weights are placed after the phoneme: <code>{`C = p*1 t*3 k`}</code>
                </li>
                <li>
                    Manually-marked weights on phonemes are carried over into any categories using that phoneme (i.e.{" "}
                    <code>{`P = p*4 k; C = $P t`}</code> will have phoneme <code>{`p`}</code> have weight 4 in all
                    categories).
                </li>
                <li>
                    Weighted optionals define what chance for that optional to appear (i.e. a weight of 33 means that it
                    will appear ~33% of the time, 1 means 1%, etc).
                </li>
                <li>Selection elements and phonemes without weights are defaulted to 1.</li>
            </ul>
            <h3>Phonology</h3>
            <ul>
                <li>
                    Category: <code>{`name = phoneme phoneme phoneme ...`}</code>
                </li>
                <li>Phonemes can be any character that is not whitespace or a semicolon </li>
                <li>
                    Using a category in a category: <code>{`C = $A $B raw $D ...`}</code>
                </li>
                <li>
                    Category reference loops are not allowed (e.g. <code>{`C = p t $K; K = a b $C`}</code>,{" "}
                    <code>{`V = a e $V o u`}</code>, etc).
                </li>
            </ul>
            <h3>Syllable</h3>
            <ul>
                <li>
                    Syllable definition: <code>{`syllable: components...`}</code>
                </li>
                <li>
                    Using categories in syllables: <code>{`syllable: $C $V $N`}</code>
                </li>
                <li>
                    Grouping components:{" "}
                    <code>{`
                        syllable: {group}
                        {$C$V}
                        {}
                    `}</code>
                </li>
                <li>
                    Optional: <code>{`syllable: (s)$C$V(n)`}</code> Defaults to a 50% chance of appearing.
                </li>
                <li>
                    Selection: <code>{`syllable: [$P,$F]r$V`}</code>
                </li>
                <li>
                    Optional with weight: <code>{`syllable: ($C)*5$V`}</code>
                </li>
                <li>
                    Selection with weight: <code>{`syllable: [$P,$F*3,{$K$L}*1]r$V`}</code>
                </li>
            </ul>
            <h3>Named Components</h3>
            <ul>
                <li>
                    Named component definition: <code>{`component: name = syllableComponents...`}</code>
                </li>
                <li>You can use any valid syllable component as a named component.</li>
                <li>
                    Use a name component by using a percent sign before the name: <code>{`%name`}</code>
                </li>
                <li>Named components can be used anywhere a syllable component can be used.</li>
                <li>
                    Named component reference loops are not allowed. For example:{" "}
                    <code>{`component: a = a %b; component: b = %a b`}</code> (two components referencing each other),
                    <code>{`component: name = %name $C`}</code> (component referencing itself).
                </li>
            </ul>
            <h3>Rejections</h3>
            <ul>
                <li>Can reject a word based on category or phoneme</li>
                <li>
                    After <code>{`reject: `}</code>, place any series of components: <code>{`reject: $C$V`}</code>
                </li>
                <li>
                    You can have multiple <code>{`reject: `}</code> directives on multiple lines
                </li>
                <li>
                    Separate multiple rejects on the same line with vertical bar <code>{`|`}</code>, surrounding each
                    with curly brackets:
                    <code>{`reject: $C$V|$V$V`}</code>
                </li>
                <li>
                    Checking for components at the start and end of words is possible:{" "}
                    <code>{`reject: ^start|end&`}</code>
                </li>
            </ul>
            {/* <h3>UNIMPLMENTED: Replacements</h3>
                <ul>
                    <li>You can replace characters or syllable components with other <b>raw</b> characters.</li>
                    <li>The form of a replacement rule is: <code>{`replace: {source} > {substitution} / condition // optionalException`}</code></li>
                    <li>The source and substitution can be any character.</li>
                    <li>The source only can be any syllable component.</li>
                    <li>Sources and substitutions can be empty; just have
                        nothing within the curly brackets.</li>
                    <li>Empty sources means that the substitution will be
                        inserted where the condition matches.</li>
                    <li>Empty substitutions means that the source will be
                        deleted.</li>
                    <li>Conditions are <i>required</i>. Exceptions are
                        <i>optional</i>.
                    </li>
                    <li>Conditions specify where and when to apply a
                        replacement.</li>
                    <li>Exceptions work similarly to conditions, in that you can
                        specify when <i>not</i> to apply a replacement.</li>
                    <li>Conditions and exceptions must have exactly one
                        underscore (<code>{`_`}</code>) representing the source.
                    </li>
                    <li>Conditions and exceptions with only an underscore means
                        that it will always run when the source matches.</li>
                    <li>Conditions and exceptions can accept syllable
                        components. For example this will replace "c" with "qu"
                        before "a", "e", and "o":
                        <code>{`replace: {c} > {qu} / _ [a,e,o,]`}</code>
                    </li>
                    <li>Conditions and exceptions can also accept word start and
                        word end checks. This will replace leading "s" with
                        "es": <code>{`replace: {s} > {es} / ^ _`}</code></li>
                </ul> */}
            <h3>Letters</h3>
            <ul>
                <li>
                    Have a line with <code>{`letters: `}</code> to define a sort order for your words.
                </li>
                <li>Each "letter" can have multiple characters.</li>
                <li>
                    Only the last <code>{`letters: `}</code> directive will be applied.
                </li>
            </ul>
        </section>
    );
}
