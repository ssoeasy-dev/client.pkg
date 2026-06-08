import { AuthManager } from "@ssoeasy-dev/core";
import { createContext } from "react";

export const AuthContext = createContext<AuthManager | null>(null);
