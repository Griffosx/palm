import { HashRouter, Routes, Route } from "react-router-dom";
import "./App.css";

// Import components
import SplashScreen from "./components/SplashScreen";
import InboxPage from "./pages/inbox/Inbox";

function App() {
  return (
    <HashRouter basename={"/"}>
      <Routes>
        <Route path="/" element={<SplashScreen />} />
        <Route path="/inbox" element={<InboxPage />} />
      </Routes>
    </HashRouter>
  );
}

export default App;
