import React, { useMemo } from "react"
import { AuthConfig, AuthManager, AuthState } from "@ssoeasy-dev/core";
import { useEffect, useState } from "react";
import { AuthContext } from "./context";

export const AuthProvider: React.FC<{
  config: AuthConfig;
  children: React.ReactNode;
}> = ({ config, children }) => {
  const auth = useMemo(() => new AuthManager(config), [config.serviceId, config.redirectUri]);
  const [state, setState] = useState<AuthState>(auth.getState());

  useEffect(() => {
    const unsubscribe = auth.onStateChange(setState);
    return unsubscribe;
  }, [auth]);

  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
};
