import { useState } from "react";
import { Route, Routes, BrowserRouter as Router } from "react-router-dom";
import { Login } from "./pages/Login";

function App() {
  const [count, setCount] = useState(0);

  return (
    <div>
      <Router>
        <Routes>
          <Route path="/*" element={<Login />} />
        </Routes>
      </Router>
    </div>
  );
}

export default App;
