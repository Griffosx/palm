import { HashRouter, Routes, Route } from "react-router-dom";
import "./App.css";

// Import components
import SplashScreen from "./components/SplashScreen";
import HomePage from "./pages/HomePage";

function App() {
  return (
    <HashRouter basename={"/"}>
      <Routes>
        <Route path="/" element={<SplashScreen />} />
        <Route path="/home" element={<HomePage />} />
      </Routes>
    </HashRouter>
  );
}

export default App;
