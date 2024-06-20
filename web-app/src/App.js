import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Navbar from "./components/Navbar";
import ProducerPage from "./pages/ProducerPage";
import ConsumerPage from "./pages/ConsumerPage";
import BrokerPage from "./pages/BrokerPage";
import LoginPage from "./pages/LoginPage";

function App() {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path="/producer" element={<ProducerPage />} />
        <Route path="/consumer" element={<ConsumerPage />} />
        <Route path="/broker" element={<BrokerPage />} />
        <Route path="/login" element={<LoginPage />} />
      </Routes>
    </Router>
  );
}

export default App;
