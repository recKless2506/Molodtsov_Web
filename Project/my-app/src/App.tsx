// src/App.tsx
import { Routes, Route } from "react-router-dom";
import { AppNavbar } from "./components/AppNavbar/AppNavbar";
import { HomePage } from "./pages/HomePage/HomePage";
import ITunesPage from "./pages/ITunesPage/ITunesPage";
import { AlbumPage } from "./pages/AlbumPage/AlbumPage";

function App() {
  return (
    <>
      <AppNavbar />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/itunes" element={<ITunesPage />} />
        <Route path="/album/:id" element={<AlbumPage />} />
      </Routes>
    </>
  );
}

export default App;
