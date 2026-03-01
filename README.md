# client.pkg

TypeScript SDK для подключения внешних приложений к SSO Easy. Монорепозиторий с фреймворк-агностичным ядром и адаптерами под конкретные фреймворки.

## Пакеты

| Пакет   | npm                  | Описание                                                   |
| ------- | -------------------- | ---------------------------------------------------------- |
| `core`  | `@ssoeasy-dev/core`  | Фреймворк-агностичный `AuthManager` — OAuth 2.0 PKCE flow  |
| `react` | `@ssoeasy-dev/react` | React-адаптер: `AuthProvider`, `useAuth`, `ProtectedRoute` |

Планируется: `@ssoeasy-dev/vue`, `@ssoeasy-dev/svelte`.

## Структура репозитория

```
client.pkg/
├── core/                  # @ssoeasy-dev/core
│   ├── src/
│   │   ├── AuthManager.ts
│   │   ├── types.ts
│   │   ├── storage.ts     # sessionStorage helpers
│   │   ├── utils.ts       # PKCE utils (verifier, challenge, state)
│   │   └── constants.generated.ts
│   └── package.json
├── react/                 # @ssoeasy-dev/react
│   ├── src/
│   │   ├── AuthProvider.tsx
│   │   ├── ProtectedRoute.tsx
│   │   ├── useAuth.ts
│   │   ├── context.ts
│   │   └── index.ts
│   └── package.json
└── examples/
    └── react/             # Пример использования с React Router
```

## Лицензия

MIT — см. [LICENSE](LICENSE).

## Контакты

- Email: morewiktor@yandex.ru
- Telegram: [@MoreWiktor](https://t.me/MoreWiktor)
- GitHub: [@MoreWiktor](https://github.com/MoreWiktor)
