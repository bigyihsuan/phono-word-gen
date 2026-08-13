import Footer from "./Footer";
import Header from "./Header";
import Input from "./Input";
import Output from "./Output";

export default function App() {
    return (
        <>
            <Header />
            <div className="main">
                <Input />
                <Output />
            </div>
            <Footer />
        </>
    );
}
