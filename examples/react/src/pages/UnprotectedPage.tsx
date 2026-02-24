import { useCallback } from "react";
import { useNavigate } from "react-router-dom";

export const UnprotectedPage = () => {
  const navigate = useNavigate();

  const toProtectedPage = useCallback(() => {
    navigate("/protected")
  }, [navigate]);

  return (
    <div>
      UnprotectedPage
      <button onClick={toProtectedPage}>
        To Protected Page
      </button>
    </div>
  );
};
