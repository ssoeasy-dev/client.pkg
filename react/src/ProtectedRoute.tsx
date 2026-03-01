import React, { useEffect, useState } from "react";
import { useAuth } from "./useAuth";

export const ProtectedRoute: React.FC<{
  children: React.ReactNode;
  redirectTo?: string;
  fallback?: React.ReactNode;
  unauthenticatedElement?: React.ReactNode;
}> = ({
  children,
  redirectTo,
  fallback = <div>Loading...</div>,
  unauthenticatedElement,
}) => {
  const auth = useAuth();
  const [isAuthorized, setIsAuthorized] = useState<boolean | null>(null);

  useEffect(() => {
    auth.checkAuth().then((val) => {
      setIsAuthorized(val);
    });
  }, [auth]);

  if (isAuthorized === null) return fallback;

  if (!isAuthorized) {
    if (unauthenticatedElement !== undefined) {
      return <>{unauthenticatedElement}</>;
    }
    auth.login({ redirectTo });
    return null;
  }

  return <>{children}</>;
};
