# IUSTITIA - Defence CTF Service

> 2187 год. На Марсе существовало 11 государств. Одиннадцатое в течение 20 лет тайно внедряло своих агентов в остальные 10 куполов под видом обычных граждан. В 2187 году началась открытая война. К 2187-му девять государств пали. Последние выжившие трибуналы ещё работают.
>
> Все внедрённые агенты Одиннадцатого были разоблачены и осуждены трибуналами выживших государств - их дела хранятся в системе **IUSTITIA v3.1**, разработанной нейтральным **Межкупольным Технологическим Бюро**. Пока они сидят - коллаборационистские ячейки внутри куполов парализованы.
>
> Одиннадцатое начало операцию **«Свободный Марс»**: взломать трибуналы противников, добраться до секретных приложений к делам агентов и изменить их статус. Освобождённые агенты активируют спящие ячейки. Купол падёт изнутри.

---

## Архитектура

### Backend

Clean Architecture, строгое разделение слоёв:

```text
HTTP (chi + oapi-codegen)
        │
   controller/restapi      <- транспорт, авторизация, валидация
        │
      usecase              <- бизнес-логика, проверки ролей
        │
       repo                <- sqlc (статические запросы) + squirrel (динамика)
        │
      SQLite               <- modernc.org/sqlite (pure-Go)
```

- **DI:** `google/wire` (compile-time), без глобального состояния.
- **Миграции:** `pressly/goose` (`migrations/001_init.sql`, `002_seed.sql`).
- **Спецификация:** OpenAPI 3 в `backend/codegen/`, кодогенерация типов и сервера через `oapi-codegen`.
- **Аутентификация:** JWT (`golang-jwt/jwt v4`), bcrypt-хэш паролей (`golang.org/x/crypto`).
- **Логирование:** `zerolog` + `lumberjack` (ротация).

### Frontend

React 18 + TypeScript + Vite, организация по **Feature-Sliced Design (FSD)**:

```text
src/app          <- bootstrap, провайдеры, роутер
src/pages        <- маршруты (страницы трибунала)
src/widgets      <- композиция блоков для страниц
src/features     <- пользовательские сценарии (login, submit-complaint, …)
src/entities     <- бизнес-сущности (case, user, document, complaint, archive)
src/shared       <- UI-кит, api-клиент, libs, конфиги
```

- **API-клиент:** `openapi-fetch` поверх типов, сгенерированных через `openapi-typescript` из `backend/internal/openapi/openapi.yml` (скрипт `npm run gen:api`).
- **Server state:** TanStack Query (`@tanstack/react-query`).
- **Client state:** Zustand.
- **Роутинг:** `react-router-dom` v6.
- **Формы и валидация:** `react-hook-form` + `zod` (через `@hookform/resolvers`).
- **Безопасность рендера:** `dompurify` для пользовательского HTML.
- **Иконки/шрифты:** `lucide-react`, `@fontsource/*`.
- **Тесты:** `vitest` + `@testing-library/react` + `jsdom`.
- **Линт/формат:** `eslint` (typescript-eslint, react, react-hooks), `prettier`, `stylelint`.
- **Прод-раздача:** статика билдится в Docker, отдаётся `nginx:1.27-alpine`, `/api/*` проксируется на `iustitia-backend:8080`.

## Структура

```text
services/iustitia/
├── docker-compose.yml         # backend :8080, frontend :8081
├── backend/
│   ├── cmd/app/               # entrypoint
│   ├── codegen/               # oapi-codegen + sqlc конфиги
│   ├── config/                # загрузка ENV
│   ├── migrations/            # goose SQL
│   ├── queries/               # sqlc *.sql
│   └── internal/
│       ├── app/               # bootstrap, http server
│       ├── apperr/            # доменные ошибки
│       ├── controller/restapi # HTTP handlers
│       ├── domain/            # User, Case, Document, Complaint, Archive
│       ├── openapi/           # сгенерированные модели
│       ├── repo/              # sqlstore + sqlc
│       ├── usecase/           # бизнес-слой
│       └── wire/              # DI
└── frontend/
    ├── Dockerfile             # multi-stage: vite build -> nginx:1.27-alpine
    ├── nginx.conf             # reverse proxy /api/* -> iustitia-backend:8080
    ├── index.html             # vite entry
    ├── vite.config.ts         # сборка
    ├── package.json           # React 18 + TS + Vite + TanStack Query + Zustand
    ├── eslint.config.js       # eslint flat config
    └── src/
        ├── app/               # bootstrap, провайдеры, роутер
        ├── pages/             # маршруты
        ├── widgets/           # композиция блоков
        ├── features/          # сценарии
        ├── entities/          # case, user, document, complaint, archive
        └── shared/            # ui-кит, api (openapi-fetch), libs
```

## Роли пользователей

| Username           | Password       | Роль       | Перевод     |
|--------------------|----------------|------------|-------------|
| `citizen_07`       | `c1t1z3n!`     | citizen    | гражданин   |
| `prosecutor_11`    | `pr0s3cut0r!`  | prosecutor | прокурор    |
| `prosecutor_12`    | `pr0s3cut0r!`  | prosecutor | прокурор    |
| `prosecutor_13`    | `pr0s3cut0r!`  | prosecutor | прокурор    |
| `judge_3`          | `ju$tice_189`  | judge      | судья       |
| `judge_4`          | `ju$tice_189`  | judge      | судья       |
| `registrar_aria7`  | `r3g1str4r!`   | registrar  | регистратор |

---

## Whitelist

Список ресурсов, которые должны быть открыты участникам Defense CTF для работы с сервисами `services/iustitia`.

### Go

- go.dev
- pkg.go.dev
- golang.org
- proxy.golang.org
- sum.golang.org
- github.com

### Backend libs

- github.com/go-chi/chi
- github.com/oapi-codegen/oapi-codegen
- github.com/oapi-codegen/runtime
- github.com/getkin/kin-openapi
- github.com/Masterminds/squirrel
- github.com/pressly/goose
- github.com/golang-jwt/jwt
- github.com/google/wire
- github.com/google/uuid
- github.com/go-playground/validator
- github.com/wahrwelt-kit/go-httpkit
- github.com/wahrwelt-kit/go-logkit
- github.com/rs/zerolog
- gopkg.in/natefinch/lumberjack.v2
- github.com/sethvargo/go-retry
- go.uber.org/multierr
- golang.org/x/crypto
- golang.org/x/sync
- modernc.org/sqlite
- sqlc.dev
- github.com/sqlc-dev/sqlc

### SQLite

- sqlite.org

### OpenAPI

- swagger.io
- spec.openapis.org
- editor.swagger.io

### Frontend libs

- npmjs.com
- registry.npmjs.org
- nodejs.org
- vitejs.dev
- react.dev
- typescriptlang.org
- tanstack.com (react-query)
- zustand-demo.pmnd.rs / docs.pmnd.rs/zustand
- react-hook-form.com
- zod.dev
- reactrouter.com
- openapi-ts.dev (openapi-typescript / openapi-fetch)
- github.com/cure53/DOMPurify
- lucide.dev
- fontsource.org
- vitest.dev
- testing-library.com
- eslint.org
- typescript-eslint.io
- prettier.io
- stylelint.io

### Docker

- hub.docker.com
- docs.docker.com
