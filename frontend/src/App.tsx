import Footer from "./Footer";
import Header from "./Header";
import Input from "./Input";
import Output from "./Output";
import QuickDocs from "./QuickDocs";

export default function App() {
    return (
        <>
            <Header />
            <main>
                <Input />
                <Output />
                <QuickDocs />
            </main>
            <Footer />
        </>
    );
}
