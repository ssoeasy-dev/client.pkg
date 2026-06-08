# @ssoeasy-dev/core

Фреймворк-агностичная библиотека для SSO-авторизации в браузере. Реализует OAuth 2.0 PKCE flow: логин с редиректом, обмен кода на токены, автоматический refresh, очередь запросов при обновлении токена.

> Только для браузерного окружения. Использует `window`, `sessionStorage` и `crypto` API.

## Установка

```bash
npm install @ssoeasy-dev/core
```

## Быстрый старт

```typescript
import { AuthManager } from "@ssoeasy-dev/core";

const auth = new AuthManager({
  serviceId: "your-service-id",
  redirectUri: "https://yourapp.com",
  loginPath: "/callback", // auth.front вернёт код на redirectUri + loginPath
  authPageUrl: "https://sso.example.com/login", // SSO страница логина
  authServerConfig: {
    baseURL: "https://api.example.com",
  },
});

// Проверка аутентификации (пытается refresh по cookie если токена нет)
const isAuth = await auth.checkAuth();
if (!isAuth) {
  auth.login({ redirectTo: "/dashboard" });
}

// На странице callback (redirectUri + loginPath)
const { redirectTo } = await auth.handleRedirectCallback();
window.location.href = redirectTo;

// Axios instance с автоматическим добавлением токена и retry при 401
const api = auth.getClient();
const response = await api.get("/api/data");

// Выход
await auth.logout();
```

## Конфигурация

```typescript
interface AuthConfig {
  serviceId: string; // Идентификатор сервиса
  redirectUri: string; // Базовый URL приложения (https://yourapp.com)
  loginPath: string; // Путь callback-страницы (/callback)
  authPageUrl?: string; // URL SSO страницы логина (default: http://localhost:5173/login)
  clientConfig?: CreateAxiosDefaults; // Конфигурация axios для клиентских запросов
  authServerConfig?: {
    baseURL?: string; // Base URL auth.api (default: https://auth.ssoeasy.ru)
    endpoints?: {
      authorize?: string; // default: /api/v1/auth/authorize/:serviceId
      logout?: string; // default: /api/v1/auth/logout
    };
  };
}
```

## API

| Метод                                                       | Описание                                                                                                         |
| ----------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `checkAuth(): Promise<boolean>`                             | Проверяет наличие токена в памяти. Если нет — пытается refresh через httpOnly cookie.                            |
| `login(options?: { redirectTo?: string }): void`            | Генерирует PKCE verifier/challenge, сохраняет в sessionStorage, redirect на SSO.                                 |
| `handleRedirectCallback(): Promise<{ redirectTo: string }>` | Читает `code` и `state` из URL, обменивает на токены, сохраняет access token в памяти.                           |
| `logout(): Promise<void>`                                   | Вызывает logout endpoint, очищает access token из памяти.                                                        |
| `getClient(): AxiosInstance`                                | Axios instance для запросов к API — автоматически добавляет `Authorization` заголовок и выполняет retry при 401. |
| `getAccessToken(): string \| null`                          | Возвращает текущий access token из памяти.                                                                       |
| `getState(): AuthState`                                     | Возвращает текущее состояние: `{ isLoading, isAuthenticated, error }`.                                           |
| `onStateChange(fn): () => void`                             | Подписка на изменения состояния. Возвращает функцию отписки.                                                     |

## Состояние

```typescript
interface AuthState {
  isLoading: boolean; // true во время refresh или начальной проверки
  isAuthenticated: boolean; // true если в памяти есть access token
  error: Error | null; // последняя ошибка
}

// Пример подписки
const unsubscribe = auth.onStateChange((state) => {
  if (state.error) console.error("Auth error:", state.error);
});
// Отписка
unsubscribe();
```

## Хранение токенов

| Токен         | Где хранится                      | Кто управляет                           |
| ------------- | --------------------------------- | --------------------------------------- |
| Access token  | Память (переменная класса)        | `AuthManager`                           |
| Refresh token | httpOnly cookie                   | Сервер (`auth.api`)                     |
| PKCE verifier | `sessionStorage` (ключ — `state`) | `AuthManager`, удаляется после callback |

При перезагрузке страницы access token теряется — `checkAuth()` автоматически пытается получить новый через refresh cookie.

## Обработка 401

`getClient()` возвращает Axios instance с перехватчиком: при получении 401 автоматически выполняется refresh, затем повторяется исходный запрос. Параллельные запросы во время refresh ставятся в очередь и выполняются после его завершения.

## PKCE Flow

```
1. login()
   → генерирует verifier + challenge (SHA-256)
   → сохраняет { verifier, redirectTo } в sessionStorage по ключу state
   → redirect → auth.front /login?challenge=...&redirect_uri=...&service_id=...&state=...

2. auth.front → auth.api → auth.svc
   → проверка credentials, сохранение auth_code
   → redirect → redirectUri + loginPath?code=...&state=...

3. handleRedirectCallback()
   → читает code + state из URL
   → берёт verifier из sessionStorage
   → POST auth.api /api/v1/auth/authorize/:serviceId { code, verifier }
   → получает access token из Authorization header
   → сохраняет в памяти
   → возвращает { redirectTo }
```

## Лицензия

MIT — см. [LICENSE](../LICENSE).

## Контакты

- Email: morewiktor@yandex.ru
- Telegram: [@MoreWiktor](https://t.me/MoreWiktor)
- GitHub: [@MoreWiktor](https://github.com/MoreWiktor)
