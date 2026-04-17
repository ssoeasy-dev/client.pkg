import type { CreateAxiosDefaults } from "axios";

export interface AuthState {
  isLoading: boolean;
  isAuthenticated: boolean;
  error: Error | null;
}

export type StateChangeListener = (state: AuthState) => void;

export interface AuthServerEndpoints {
  authorize?: string;
  logout?: string;
  me?: string;
}

interface AuthServerConfig {
  baseURL?: string;
  endpoints?: AuthServerEndpoints;
}

export interface AuthConfig {
  /** Идентификатор сервиса (клиента) */
  serviceId: string;
  /** Полный URL, на который вернётся SSO после логина */
  redirectUri: string;
  loginPath: string;
  /** URL страницы входа (по умолчанию http://localhost:5173/login) */
  authPageUrl?: string;
  clientConfig?: CreateAxiosDefaults;
  authServerConfig?: AuthServerConfig;
}

export interface MeResponse {
  id: string;
  firstname: string;
  lastname: string;
  login: string;
  companies: {
    id: string;
    name: string;
    serviceName?: string;
    subscriptionId?: string;
    subscriptionIsActive?: boolean;
  }[];
}
