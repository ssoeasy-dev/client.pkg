# client.pkg

TypeScript SDK для подключения внешних приложений к SSO Easy. Монорепозиторий с фреймворк-агностичным ядром и адаптерами под конкретные фреймворки.

## Пакеты

| Пакет          | npm                  | Описание                                                                |
| -------------- | -------------------- | ----------------------------------------------------------------------- |
| `client/core`  | `@ssoeasy-dev/core`  | Фреймворк-агностичный `AuthManager` — OAuth 2.0 PKCE flow для фронтенда |
| `client/react` | `@ssoeasy-dev/react` | React-адаптер: `AuthProvider`, `useAuth`, `ProtectedRoute`              |
| `api/go/core`  | `-`                  | Client и middleware для проверки прав доступа на сервере                |
| `api/go/echo`  | `-`                  | Адаптер для фреймворка echo                                             |

## Структура репозитория

```
client.pkg/
├── client/                    # Пакеты для клиентской части
│   ├── core/                   # Ядро функционала npm:@ssoeasy-dev/core
│   ├── examples/               # Примеры использования
│   └── react/                  # Адаптер для react npm:@ssoeasy-dev/react
└── api/                       # Пакеты для серверной части
    ├── examples/               # Примеры использования
    └── go/                     # Пакеты для go
        ├── core/                # Ядро функционала
        └── echo/                # Адаптер для echo
```

## Лицензия

MIT — см. [LICENSE](LICENSE).

## Контакты

- Email: morewiktor@yandex.ru
- Telegram: [@MoreWiktor](https://t.me/MoreWiktor)
- GitHub: [@MoreWiktor](https://github.com/MoreWiktor)
