import { ProtectedRoute } from "@ssoeasy-dev/react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { LoginPage } from "./pages/LoginPage";
import { ProtectedPage } from "./pages/ProtectedPage";
import { UnprotectedPage } from "./pages/UnprotectedPage";

export const AppRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={<UnprotectedPage/>} />
        <Route
          path="/protected"
          element={
            <ProtectedRoute>
              <ProtectedPage/>
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
};
