import Footer from "./Footer";
import Header from "./Header";
import Input from "./Input";
import Output from "./Output";
import QuickReference from "./QuickDocs";

export default function App() {
    return (
        <>
            <Header />
            <main className="container">
                <div className="row">
                    <Input />
                    <Output output={""} generatedCount={-1} duplicateCount={-1} rejectedCount={-1} replacedCount={-1} />
                </div>
                <hr></hr>
                <QuickReference />
            </main>
            <Footer />
        </>
    );
}
