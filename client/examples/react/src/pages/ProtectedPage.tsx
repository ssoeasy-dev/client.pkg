import { useCallback, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@ssoeasy-dev/react";

export const ProtectedPage = () => {
  const auth = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const logout = useCallback(() => {
    setLoading(true);
    auth
      .logout()
      .catch((error) => console.error(error))
      .finally(() => {
        setLoading(false);
        navigate("/");
      });
  }, [auth, navigate]);

  const toUnrotectedPage = useCallback(() => {
    navigate("/")
  }, [navigate]);

  return (
    <div>
      ProtectedPage
      <button onClick={logout} disabled={loading}>
        {loading ? "Loading..." : "Logout"}
      </button>
      <button onClick={toUnrotectedPage}>
        To Unrotected Page
      </button>
    </div>
  );
};
