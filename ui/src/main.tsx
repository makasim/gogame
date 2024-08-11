import { createRoot } from "react-dom/client";
import { Navigate, Route, HashRouter as Router, Routes } from "react-router-dom";

import "./index.scss";

import { App } from "./App.tsx";
import { UserForm } from "./UserForm.tsx";
import { GamesPage } from "./GamesPage.tsx";

createRoot(document.getElementById("root")!).render(
  <Router>
    <Routes>
      <Route path="/" element={<UserForm />} />
      <Route path="player/:playerId" element={<GamesPage />} />
      <Route path="player/:playerId/game/:gameId" element={<App />} />
      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  </Router>,
);
