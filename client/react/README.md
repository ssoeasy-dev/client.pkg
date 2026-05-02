# @ssoeasy-dev/react

React-адаптер для SSO Easy. Оборачивает `@ssoeasy-dev/core` в React Context и предоставляет хук `useAuth` и компонент `ProtectedRoute`.

## Установка

```bash
npm install @ssoeasy-dev/react @ssoeasy-dev/core
```

## Быстрый старт

### 1. Оборачиваем приложение в `AuthProvider`

```tsx
import { AuthProvider } from "@ssoeasy-dev/react";
import { AuthManager } from "@ssoeasy-dev/core";

const auth = new AuthManager({
  serviceId: "your-service-id",
  redirectUri: "https://yourapp.com",
  loginPath: "/callback",
  authPageUrl: "https://sso.example.com/login",
  authServerConfig: {
    baseURL: "https://api.example.com",
  },
});

function App() {
  return (
    <AuthProvider auth={auth}>
      <AppRouter />
    </AuthProvider>
  );
}
```

### 2. Защищаем роуты через `ProtectedRoute`

```tsx
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ProtectedRoute } from "@ssoeasy-dev/react";

function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/callback" element={<CallbackPage />} />
        <Route path="/" element={<PublicPage />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute redirectTo="/dashboard">
              <Dashboard />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}
```

### 3. Используем `useAuth` в компонентах

```tsx
import { useAuth } from "@ssoeasy-dev/react";

function Header() {
  const auth = useAuth();

  return (
    <header>
      {auth.getState().isAuthenticated ? (
        <button onClick={() => auth.logout()}>Выйти</button>
      ) : (
        <button onClick={() => auth.login()}>Войти</button>
      )}
    </header>
  );
}
```

### 4. Страница callback

```tsx
import { useEffect } from "react";
import { useAuth } from "@ssoeasy-dev/react";
import { useNavigate } from "react-router-dom";

function CallbackPage() {
  const auth = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    auth.handleRedirectCallback().then(({ redirectTo }) => {
      navigate(redirectTo, { replace: true });
    });
  }, []);

  return <div>Авторизация...</div>;
}
```

## API

### `AuthProvider`

```tsx
<AuthProvider auth={AuthManager}>{children}</AuthProvider>
```

Предоставляет `AuthManager` через React Context. Должен оборачивать всё приложение.

### `useAuth`

```tsx
const auth = useAuth(); // возвращает AuthManager
```

Хук для доступа к `AuthManager` внутри компонентов. Выбрасывает ошибку если вызван вне `AuthProvider`.

### `ProtectedRoute`

```tsx
<ProtectedRoute
  redirectTo="/target-path" // путь для redirectTo при логине (опционально)
  fallback={<Spinner />} // отображается во время checkAuth (по умолчанию <div>Loading...</div>)
  unauthenticatedElement={<Node />} // рендерится вместо редиректа на SSO, если пользователь не авторизован (опционально)
>
  <ProtectedContent />
</ProtectedRoute>
```

При монтировании вызывает `auth.checkAuth()`. Пока идёт проверка — рендерит `fallback`.

Если пользователь не авторизован:

- **`unauthenticatedElement` не передан** — вызывает `auth.login({ redirectTo })`, редирект на SSO (поведение по умолчанию)
- **`unauthenticatedElement` передан** — рендерит переданный элемент без редиректа

#### Пример: внутренняя страница выбора входа/регистрации

Если в приложении есть собственная страница `/auth` с кнопками «Войти» и «Зарегистрироваться», можно передать `<Navigate>` вместо того, чтобы сразу уходить на SSO:

```tsx
import { Navigate } from "react-router-dom";
import { ProtectedRoute } from "@ssoeasy-dev/react";

const LayoutWrapper = () => {
  return (
    <ProtectedRoute unauthenticatedElement={<Navigate to="/auth" replace />}>
      <Layout>
        <Outlet />
      </Layout>
    </ProtectedRoute>
  );
};
```

Страница `/auth` при этом остаётся публичной — она не обёрнута в `ProtectedRoute`.

## Экспорты

```typescript
export { AuthProvider } from "./AuthProvider";
export { ProtectedRoute } from "./ProtectedRoute";
export { useAuth } from "./useAuth";
```

## Лицензия

MIT — см. [LICENSE](../LICENSE).

## Контакты

- Email: morewiktor@yandex.ru
- Telegram: [@MoreWiktor](https://t.me/MoreWiktor)
- GitHub: [@MoreWiktor](https://github.com/MoreWiktor)
