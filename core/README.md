# @ssoeasy-dev/core

**Core library** for SSO authentication in browser using OAuth 2.0 with PKCE.  
Provides a framework-agnostic `AuthManager` class that handles the full authentication flow:

- Login with redirect to SSO server
- Token exchange (authorization code + PKCE verifier)
- Automatic token refresh using refresh token (httpOnly cookie)
- Request queuing during token refresh
- Logout

Works with **React**, **Vue**, **Svelte** (adapters available separately).

---

## Installation

```bash
npm install @ssoeasy-dev/core
# or
yarn add @ssoeasy-dev/core
```

> **Note:** This library is intended for browser environments only. It uses `window`, `sessionStorage`, and `crypto` APIs.

---

## Quick Start

```typescript
import { AuthManager } from '@ssoeasy-dev/core';

const auth = new AuthManager({
  serviceId: 'your-client-id',
  redirectUri: 'https://yourapp.com/callback',
  baseURL: 'https://api.example.com',
  loginUrl: 'https://sso.example.com/login', // optional, default: http://localhost:5173/login
});

// On your login page / protected route
auth.checkAuth().then(isAuthenticated => {
  if (!isAuthenticated) {
    auth.login({ redirectTo: '/dashboard' });
  }
});

// On your callback page (redirectUri)
try {
  const { redirectTo } = await auth.handleRedirectCallback();
  window.location.href = redirectTo;
} catch (err) {
  console.error('Authentication failed', err);
}

// Logout
await auth.logout();

// Use pre-configured axios instance with automatic token injection and refresh
const api = auth.getAxiosInstance();
api.get('/user').then(response => console.log(response.data));
```

---

## API Reference

### `AuthManager`

#### Constructor

```typescript
new AuthManager(config: AuthConfig)
```

**AuthConfig**

| Property      | Type                               | Description                                                                 |
|---------------|------------------------------------|-----------------------------------------------------------------------------|
| `serviceId`   | `string`                           | Your service/client identifier.                                             |
| `redirectUri` | `string`                           | Full URL where the SSO server will redirect back after login.               |
| `loginUrl`    | `string` (optional)                | SSO login page URL. Default: `http://localhost:5173/login`.                 |
| `baseURL`     | `string` (optional)                | Base URL for API requests. Used for token refresh and logout.               |
| `endpoints`   | `{ authorize?: string; logout?: string }` (optional) | Custom API endpoint paths. Defaults: `/api/v1/auth/authorize`, `/api/v1/auth/logout`. |

---

#### Methods

| Method                     | Description                                                                                          |
|----------------------------|------------------------------------------------------------------------------------------------------|
| `checkAuth(): Promise<boolean>` | Checks if the user is authenticated. Attempts to refresh tokens if possible. Updates internal state. |
| `login(options?: { redirectTo?: string }): Promise<void>` | Redirects to SSO login page. Saves PKCE verifier and `redirectTo` in sessionStorage.                 |
| `handleRedirectCallback(): Promise<{ redirectTo: string }>` | Exchanges authorization code for tokens. Extracts access token from `Authorization` header.          |
| `logout(): Promise<void>`  | Calls logout endpoint and clears local access token.                                                 |
| `getAxiosInstance(): AxiosInstance` | Returns an Axios instance with automatic token injection and 401 handling (token refresh + retry).   |
| `getAccessToken(): string \| null` | Returns the current in-memory access token.                                                          |
| `getState(): AuthState`    | Returns current authentication state: `{ isLoading, isAuthenticated, error }`.                       |
| `onStateChange(listener: (state: AuthState) => void): () => void` | Subscribe to state changes. Returns unsubscribe function.                                            |

---

#### State

```typescript
interface AuthState {
  isLoading: boolean;        // True during token refresh or initial check
  isAuthenticated: boolean;  // True if a valid access token exists
  error: Error | null;       // Last error that occurred
}
```

---

## Framework Integration

While `@ssoeasy-dev/core` can be used directly in any framework, we provide dedicated adapters for a smoother experience:

- **React:** [`@ssoeasy-dev/react`(https://www.npmjs.com/package/@ssoeasy-dev/react) – `AuthProvider`, `useAuth`, `ProtectedRoute`]
- **Vue:** [`@ssoeasy-dev/vue`(https://www.npmjs.com/package/@ssoeasy-dev/vue) – composable `useAuth`, component `ProtectedRoute`]
- **Svelte:** [`@ssoeasy-dev/svelte`(https://www.npmjs.com/package/@ssoeasy-dev/svelte) – stores and `ProtectedRoute` component]

See the respective package READMEs for usage examples.

---

## How It Works (Flow)

1. **Initialization:**  
   Create an `AuthManager` instance with your service ID and redirect URI.

2. **Protected Route:**  
   Call `auth.checkAuth()`.  
   - If a valid access token is present in memory → **authenticated**.  
   - If not, it attempts a silent refresh using the httpOnly refresh cookie (endpoint [1]).  
   - If refresh fails, redirect to login (`auth.login()`) is triggered.

3. **Login Redirect:**  
   `auth.login()` generates a PKCE verifier and challenge, saves the verifier in `sessionStorage` (keyed by a random `state`), and redirects to the SSO login page with `challenge`, `redirect_uri`, `service_id`, and `state`.

4. **Callback Handling:**  
   On your `redirectUri`, call `auth.handleRedirectCallback()`. It reads `code` and `state` from the URL, retrieves the verifier from `sessionStorage`, exchanges the code for tokens (endpoint [3]), and stores the access token in memory. Returns the original `redirectTo` path.

5. **Authenticated Requests:**  
   Use `auth.getAxiosInstance()` for API calls. It automatically attaches the access token to requests. If a 401 occurs, it attempts to refresh the token and retries the failed request. Other requests are queued during refresh.

6. **Logout:**  
   `auth.logout()` calls the logout endpoint (endpoint [4]) and clears the local access token.

---

## Important Notes

- **Access Token Storage:** The access token is kept **only in memory**. On page reload, `isAuthenticated` becomes `false`, but `checkAuth()` will attempt a refresh using the httpOnly cookie.
- **Refresh Token:** Stored in an **httpOnly cookie** by the server. The library never reads it directly.
- **Session Storage:** Used temporarily during login to store the PKCE verifier and redirect target. It is automatically cleaned up after callback handling.
- **CORS & Credentials:** All token endpoints require `withCredentials: true`. The library sets this automatically.
- **SSR Safety:** Methods that rely on browser APIs (`login`, `handleRedirectCallback`, `checkAuth`) will throw an error if called on the server. Use them only in client-side code (e.g., inside `useEffect` in React, `onMounted` in Vue, or `onMount` in Svelte).

---

## Error Handling

The library catches most errors and updates the `error` field in the state. You can subscribe to state changes to display error messages:

```typescript
auth.onStateChange((state) => {
  if (state.error) {
    console.error('Auth error:', state.error);
    // show notification
  }
});
```

For critical failures (e.g., missing `code` or `state` in callback URL), methods like `handleRedirectCallback()` will throw synchronously.

---

## Example: Using Core Directly in a React App

```tsx
// auth.ts
import { AuthManager } from '@ssoeasy-dev/core';

export const auth = new AuthManager({
  serviceId: 'my-app',
  redirectUri: 'https://myapp.com/callback',
  baseURL: 'https://api.myapp.com',
});

// App.tsx
import { useEffect, useState } from 'react';
import { auth } from './auth';

function App() {
  const [isAuthenticated, setIsAuthenticated` = useState(false);

  useEffect(() => {
    auth.checkAuth().then(setIsAuthenticated);
  }, [`);

  if (!isAuthenticated) {
    auth.login({ redirectTo: window.location.pathname });
    return <div>Redirecting to login...</div>;
  }

  return <Dashboard />;
}

// Dashboard.tsx
const api = auth.getAxiosInstance();
api.get('/user').then(...);
```

---

## License

MIT

---

## Contributing

Issues and pull requests are welcome. Please ensure that your code passes the existing tests and is properly formatted.

---

## Support

If you encounter any problems, please `open an issue`(https://github.com/your-org/ssoeasy-dev/issues).