import Footer from "./Footer";
import Header from "./Header";
import Input from "./Input";
import Output from "./Output";
import QuickDocs from "./QuickDocs";

export default function App() {
    return (
        <>
            <Header />
            <main className="container">
                <div className="row">
                    <Input />
                    <Output output={""} generatedCount={0} duplicateCount={0} rejectedCount={0} replacedCount={0} />
                </div>
                <QuickDocs />
            </main>
            <Footer />
        </>
    );
}
