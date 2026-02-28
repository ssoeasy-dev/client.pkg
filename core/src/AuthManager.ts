import axios, {
  AxiosInstance,
  InternalAxiosRequestConfig,
  AxiosError,
} from "axios";
import {
  generateCodeVerifier,
  generateCodeChallenge,
  generateState,
} from "./utils";
import { saveTempData, getTempData, removeTempData } from "./storage";
import type {
  AuthState,
  StateChangeListener,
  AuthConfig,
  AuthServerEndpoints,
} from "./types";
import {
  AUTH_PAGE,
  AUTHORIZE_PATH,
  BASE_URL,
  LOGOUT_PATH,
} from "./constants.generated";

export class AuthManager {
  private serviceId: string;
  private redirectUri: string;
  private authPageUrl: string;
  private baseURL: string;
  private loginPath: string;
  private endpoints: Required<AuthServerEndpoints>;
  private accessToken: string | null = null;
  private listeners: StateChangeListener[] = [];
  private state: AuthState = {
    isLoading: false,
    isAuthenticated: false,
    error: null,
  };
  private authClient: AxiosInstance;
  private client: AxiosInstance;
  private refreshPromise: Promise<string | null> | null = null;
  private pendingRequests: Array<{
    resolve: (value: any) => void;
    reject: (reason?: any) => void;
    config: InternalAxiosRequestConfig;
  }> = [];

  constructor(config: AuthConfig) {
    this.serviceId = config.serviceId;
    this.redirectUri = config.redirectUri;
    this.loginPath = config.loginPath;
    this.authPageUrl = config.authPageUrl || AUTH_PAGE;
    this.baseURL = config.authServerConfig?.baseURL || BASE_URL;
    this.endpoints = {
      authorize:
        config.authServerConfig?.endpoints?.authorize || AUTHORIZE_PATH,
      logout: config.authServerConfig?.endpoints?.logout || LOGOUT_PATH,
    };

    if (typeof window === "undefined") {
      throw new Error("AuthManager constructor: window is undefined");
    }

    if (!this.redirectUri.startsWith("http")) {
      throw new Error(
        "AuthManager constructor: redirectUri must be URI, example: https://domain.com/",
      );
    }

    this.authClient = axios.create({ baseURL: this.baseURL });
    this.setupInterceptors(this.authClient);
    this.client = axios.create(config.clientConfig);
    this.setupInterceptors(this.client);
  }

  private setState(newState: Partial<AuthState>) {
    this.state = {
      ...this.state,
      ...newState,
    };
    this.listeners.forEach((fn) => fn(this.state));
  }

  public onStateChange(listener: StateChangeListener): () => void {
    this.listeners.push(listener);
    return () => {
      this.listeners = this.listeners.filter((l) => l !== listener);
    };
  }

  public getState(): AuthState {
    return { ...this.state };
  }

  private setAccessToken(token: string | null) {
    this.accessToken = token;
    this.setState({ isAuthenticated: token !== null });
  }

  public getAccessToken(): string | null {
    return this.accessToken;
  }

  private isBrowser(): boolean {
    return typeof window !== "undefined";
  }

  public async checkAuth(): Promise<boolean> {
    if (!this.isBrowser()) {
      return false;
    }

    if (this.refreshPromise) {
      const token = await this.refreshPromise;
      return token !== null;
    }

    this.setState({ isLoading: true, error: null });
    try {
      const token = await this.performRefresh();
      this.setAccessToken(token);
      return token !== null;
    } catch (err) {
      this.setAccessToken(null);
      this.setState({ error: err as Error });
      return false;
    } finally {
      this.setState({ isLoading: false });
    }
  }

  public async login(options?: { redirectTo?: string }): Promise<void> {
    if (!this.isBrowser()) {
      throw new Error("login can only be called in browser environment");
    }

    const verifier = await generateCodeVerifier();
    const challenge = await generateCodeChallenge(verifier);
    const state = generateState();
    const redirectTo =
      options?.redirectTo || window.location.pathname + window.location.search;

    saveTempData(state, { verifier, redirectTo });

    const loginUrlObj = new URL(this.authPageUrl);
    loginUrlObj.searchParams.set("challenge", challenge);
    loginUrlObj.searchParams.set(
      "redirect_uri",
      this.redirectUri + this.loginPath,
    );
    loginUrlObj.searchParams.set("service_id", this.serviceId);
    loginUrlObj.searchParams.set("state", state);

    window.location.href = loginUrlObj.toString();
  }

  public async handleRedirectCallback(): Promise<{ redirectTo: string }> {
    if (!this.isBrowser()) {
      throw new Error(
        "handleRedirectCallback can only be called in browser environment",
      );
    }

    const url = new URL(window.location.href);
    const code = url.searchParams.get("code");
    const state = url.searchParams.get("state");
    const traceId = url.searchParams.get("traceId");

    if (!code || !state) {
      throw new Error("Missing code or state in URL");
    }

    const tempData = getTempData(state);
    if (!tempData) {
      throw new Error("Invalid state or expired session");
    }

    removeTempData(state);

    try {
      const response = await this.authClient.post(
        this.baseURL +
          this.endpoints.authorize.replace(":serviceId", this.serviceId),
        {
          verifier: tempData.verifier,
          code,
        },
        {
          withCredentials: true,
          ...(traceId && { headers: { "x-trace-id": traceId } }),
        },
      );

      sessionStorage.removeItem('x-trace-id')

      const authHeader =
        response.headers["authorization"] || response.headers["Authorization"];
      if (!authHeader || typeof authHeader !== "string") {
        throw new Error("No Authorization header in response");
      }
      const match = authHeader.match(/^Bearer\s+(.+)$/i);
      if (!match) {
        throw new Error("Invalid Authorization header format");
      }
      this.setAccessToken(match[1]);

      return { redirectTo: tempData.redirectTo };
    } catch (err) {
      this.setAccessToken(null);
      throw err;
    }
  }

  public async logout(): Promise<void> {
    if (!this.isBrowser()) {
      return;
    }

    try {
      await this.authClient.post(
        this.baseURL + this.endpoints.logout,
        {},
        { withCredentials: true },
      );
    } finally {
      this.setAccessToken(null);
    }
  }

  private async performRefresh(): Promise<string | null> {
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    this.refreshPromise = (async () => {
      try {
        const response = await this.authClient.post(
          this.baseURL +
            this.endpoints.authorize.replace(":serviceId", this.serviceId),
          {},
          { withCredentials: true },
        );

        const authHeader =
          response.headers["authorization"] ||
          response.headers["Authorization"];

        if (!authHeader || typeof authHeader !== "string") {
          return null;
        }
        const match = authHeader.match(/^Bearer\s+(.+)$/i);

        return match ? match[1] : null;
      } catch (err) {
        if (axios.isAxiosError(err) && err.response?.status === 401) {
          return null;
        }
        throw err;
      } finally {
        this.refreshPromise = null;
      }
    })();

    return this.refreshPromise;
  }

  public getClient(): AxiosInstance {
    return this.client;
  }

  private setupInterceptors(instanse: AxiosInstance) {
    instanse.interceptors.request.use((config) => {
      if (this.accessToken) {
        config.headers.Authorization = `Bearer ${this.accessToken}`;
      }
      return config;
    });

    instanse.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalConfig = error.config as InternalAxiosRequestConfig & {
          _retry?: boolean;
        };
        if (
          error.response?.status !== 401 ||
          originalConfig._retry ||
          this.isRefreshRequest(originalConfig)
        ) {
          return Promise.reject(error);
        }

        originalConfig._retry = true;

        return new Promise((resolve, reject) => {
          this.pendingRequests.push({
            resolve,
            reject,
            config: originalConfig,
          });

          if (!this.refreshPromise) {
            this.performRefresh()
              .then((newToken) => {
                if (newToken) {
                  this.setAccessToken(newToken);
                  this.pendingRequests.forEach(
                    ({ resolve, reject, config }) => {
                      if (this.accessToken) {
                        config.headers.Authorization = `Bearer ${this.accessToken}`;
                      }
                      instanse(config).then(resolve).catch(reject);
                    },
                  );
                } else {
                  this.pendingRequests.forEach(({ reject }) => reject(error));
                  this.setAccessToken(null);
                }
                this.pendingRequests = [];
              })
              .catch((refreshError) => {
                this.pendingRequests.forEach(({ reject }) =>
                  reject(refreshError),
                );
                this.pendingRequests = [];
                this.setAccessToken(null);
              });
          }
        });
      },
    );
  }

  private isRefreshRequest(config: InternalAxiosRequestConfig): boolean {
    const url = config.url?.split("?")[0];
    const target =
      this.baseURL +
      this.endpoints.authorize.replace(":serviceId", this.serviceId);
    return url === target && config.method?.toLowerCase() === "post";
  }
}
