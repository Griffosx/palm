import { HashRouter, Routes, Route } from "react-router-dom";
import "./App.css";

// Import components
import SplashScreen from "./components/SplashScreen";
import HomePage from "./pages/HomePage";
import Page1 from "./pages/Page1";
import Page2 from "./pages/Page2";

// Create simple placeholder pages for menu items
const InboxPage = () => <Page1 />;

const SentPage = () => <Page2 />;

const DraftPage = () => <Page1 />;

function App() {
  return (
    <HashRouter basename={"/"}>
      <Routes>
        <Route path="/" element={<SplashScreen />} />
        <Route path="/home" element={<HomePage />} />
        <Route path="/page1" element={<Page1 />} />
        <Route path="/page2" element={<Page2 />} />

        {/* New menu item routes */}
        <Route path="/inbox" element={<InboxPage />} />
        <Route path="/sent" element={<SentPage />} />
        <Route path="/draft" element={<DraftPage />} />
      </Routes>
    </HashRouter>
  );
}

export default App;
