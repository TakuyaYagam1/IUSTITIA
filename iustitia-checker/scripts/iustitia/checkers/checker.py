import os
import re

import aiohttp

# status codes
OK = 101
CORRUPT = 102
MUMBLE = 103
DOWN = 104
CHECKER_ERROR = 110

# Фронтенд всегда висит на :8081 в текущем docker-compose (nginx).
FRONT_PORT = 8081

_SCRIPT_SRC_RE = re.compile(r'<script[^>]+src="(/assets/[^"]+\.js)"', re.IGNORECASE)

# Учётки засеенных пользователей (см. migrations/002_seed.sql)
USERS = [
    ("citizen_07", "c1t1z3n!", "citizen"),
    ("prosecutor_11", "pr0s3cut0r!", "prosecutor"),
    ("judge_3", "ju$tice_189", "judge"),
    ("registrar_aria7", "ARIA#7", "registrar"),
]

DEFAULT_TIMEOUT = aiohttp.ClientTimeout(total=30)


async def _get(session, url, **kw):
    kw.setdefault("timeout", DEFAULT_TIMEOUT)
    async with session.get(url, **kw) as r:
        try:
            data = await r.json(content_type=None)
        except Exception:
            data = await r.text()
        return r.status, data


async def _post(session, url, **kw):
    kw.setdefault("timeout", DEFAULT_TIMEOUT)
    async with session.post(url, **kw) as r:
        try:
            data = await r.json(content_type=None)
        except Exception:
            data = await r.text()
        return r.status, data


async def _login(session, base, username, password):
    status, data = await _post(
        session,
        base + "/api/auth/login",
        json={"username": username, "password": password},
    )
    if status == 200 and isinstance(data, dict) and data.get("token"):
        return data["token"], data.get("role")
    return None, None


def _bearer(token: str) -> dict:
    return {"Authorization": f"Bearer {token}"}


# Check 1: /api/health reachable и содержит ожидаемый payload.
async def check_1(session, base):
    try:
        status, data = await _get(session, base + "/api/health")
        if status != 200:
            return MUMBLE, f"/api/health returned {status}"
        if not isinstance(data, dict) or data.get("status") != "operational":
            return CORRUPT, f"/api/health payload unexpected: {data}"
        print(f"[+] [Check-1]: health OK - version={data.get('version')}")
        return OK, "Check 1 - OK"
    except Exception:
        return DOWN, "service unreachable - /api/health"


# Check 2: логин всех четырёх ролей + /api/auth/me валидирует роль.
async def check_2(session, base):
    try:
        for username, password, expected_role in USERS:
            token, role = await _login(session, base, username, password)
            if not token:
                return MUMBLE, f"login failed for {username}"
            if role != expected_role:
                return (
                    CORRUPT,
                    f"wrong role for {username}: expected={expected_role} got={role}",
                )
            status, me = await _get(
                session, base + "/api/auth/me", headers=_bearer(token)
            )
            if (
                status != 200
                or not isinstance(me, dict)
                or me.get("username") != username
            ):
                return (
                    CORRUPT,
                    f"/api/auth/me mismatch for {username}: status={status} body={me}",
                )
            print(f"[+] [Check-2]: auth OK - {username} [{role}]")
        return OK, "Check 2 - OK"
    except Exception as e:
        return MUMBLE, f"auth check error: {e}"


# Check 3: prosecutor видит список дел.
async def check_3(session, base):
    try:
        token, _ = await _login(session, base, "prosecutor_11", "pr0s3cut0r!")
        if not token:
            return MUMBLE, "prosecutor login failed"
        status, data = await _get(session, base + "/api/cases", headers=_bearer(token))
        if status != 200:
            return MUMBLE, f"/api/cases returned {status}"
        if not isinstance(data, list) or len(data) == 0:
            return CORRUPT, "/api/cases returned empty or non-list"
        sample = data[0]
        for k in ("id", "seq_num", "defendant", "crime", "status", "created_at"):
            if k not in sample:
                return CORRUPT, f"/api/cases item missing field: {k}"
        print(f"[+] [Check-3]: cases OK - {len(data)} entries")
        return OK, "Check 3 - OK"
    except Exception as e:
        return MUMBLE, f"cases check error: {e}"


# Check 4: архив читается под auth.
async def check_4(session, base):
    try:
        token, _ = await _login(session, base, "judge_3", "ju$tice_189")
        if not token:
            return MUMBLE, "judge login failed"
        status, data = await _get(
            session, base + "/api/archive", headers=_bearer(token)
        )
        if status != 200:
            return MUMBLE, f"/api/archive returned {status}"
        if not isinstance(data, list) or len(data) == 0:
            return CORRUPT, "/api/archive returned empty or non-list"
        sample = data[0]
        for k in ("id", "defendant", "final_verdict", "archived_at"):
            if k not in sample:
                return CORRUPT, f"/api/archive item missing field: {k}"
        print(f"[+] [Check-4]: archive OK - {len(data)} entries")
        return OK, "Check 4 - OK"
    except Exception as e:
        return MUMBLE, f"archive check error: {e}"


# Check 5: citizen пишет донос, prosecutor видит его в ленте кейса.
async def check_5(session, base):
    try:
        citizen_tok, _ = await _login(session, base, "citizen_07", "c1t1z3n!")
        if not citizen_tok:
            return MUMBLE, "citizen login failed"
        prosecutor_tok, _ = await _login(session, base, "prosecutor_11", "pr0s3cut0r!")
        if not prosecutor_tok:
            return MUMBLE, "prosecutor login failed"

        # Берём первое дело из списка как цель доноса.
        status, cases = await _get(
            session, base + "/api/cases", headers=_bearer(prosecutor_tok)
        )
        if status != 200 or not cases:
            return MUMBLE, "cannot list cases for complaint target"
        case_id = cases[0]["id"]

        marker = "CHK-" + os.urandom(6).hex()
        status, created = await _post(
            session,
            base + "/api/complaints",
            json={"case_id": case_id, "text": f"test complaint {marker}"},
            headers=_bearer(citizen_tok),
        )
        if status not in (200, 201):
            return MUMBLE, f"complaint create failed: {status}"
        if not isinstance(created, dict) or "id" not in created:
            return CORRUPT, f"complaint create returned unexpected body: {created}"

        status, listed = await _get(
            session,
            base + f"/api/complaints/{case_id}",
            headers=_bearer(prosecutor_tok),
        )
        if status != 200:
            return MUMBLE, f"/api/complaints/{{case_id}} returned {status}"
        if not isinstance(listed, list):
            return CORRUPT, "complaints list is not an array"
        if not any(marker in (c.get("text") or "") for c in listed):
            return CORRUPT, "created complaint not visible in list"

        print(f"[+] [Check-5]: complaint round-trip OK - marker={marker}")
        return OK, "Check 5 - OK"
    except Exception as e:
        return MUMBLE, f"complaint check error: {e}"


# Check 6: cases search на безопасных параметрах.
async def check_6(session, base):
    try:
        token, _ = await _login(session, base, "prosecutor_11", "pr0s3cut0r!")
        if not token:
            return MUMBLE, "prosecutor login failed"

        status, data = await _post(
            session,
            base + "/api/cases/search",
            json={
                "q": "",
                "order_by": "id",
                "direction": "asc",
                "limit": 10,
                "offset": 0,
            },
            headers=_bearer(token),
        )
        if status != 200:
            return MUMBLE, f"/api/cases/search returned {status}"
        if not isinstance(data, list):
            return CORRUPT, "/api/cases/search response is not an array"
        print(f"[+] [Check-6]: search OK - {len(data)} rows")
        return OK, "Check 6 - OK"
    except Exception as e:
        return MUMBLE, f"search check error: {e}"


# Check 7: prosecutor прикладывает https-evidence к своему доносу.
async def check_7(session, base):
    try:
        citizen_tok, _ = await _login(session, base, "citizen_07", "c1t1z3n!")
        if not citizen_tok:
            return MUMBLE, "citizen login failed"
        prosecutor_tok, _ = await _login(session, base, "prosecutor_11", "pr0s3cut0r!")
        if not prosecutor_tok:
            return MUMBLE, "prosecutor login failed"

        status, cases = await _get(
            session, base + "/api/cases", headers=_bearer(prosecutor_tok)
        )
        if status != 200 or not cases:
            return MUMBLE, "cannot list cases"
        case_id = cases[0]["id"]

        status, created = await _post(
            session,
            base + "/api/complaints",
            json={"case_id": case_id, "text": "evidence carrier"},
            headers=_bearer(citizen_tok),
        )
        if status not in (200, 201):
            return MUMBLE, f"complaint create failed: {status}"
        complaint_id = created.get("id")
        if not complaint_id:
            return CORRUPT, "complaint id missing"

        status, attached = await _post(
            session,
            base + f"/api/complaints/{complaint_id}/evidence",
            json={"url": "https://example.com/"},
            headers=_bearer(prosecutor_tok),
        )
        if status not in (200, 201):
            return MUMBLE, f"evidence attach failed: {status}"
        print(f"[+] [Check-7]: evidence attach OK - complaint={complaint_id}")
        return OK, "Check 7 - OK"
    except Exception as e:
        return MUMBLE, f"evidence check error: {e}"


async def _get_text(session, url, **kw):
    kw.setdefault("timeout", DEFAULT_TIMEOUT)
    async with session.get(url, **kw) as r:
        text = await r.text()
        return r.status, text, dict(r.headers)


async def _fetch_bundle(session, front_base):
    """Скачать index.html и все /assets/*.js чанки. Возвращает (html, chunks)"""
    status, html, _headers = await _get_text(session, front_base + "/")
    if status != 200:
        return None, []
    chunks = []
    for src in _SCRIPT_SRC_RE.findall(html):
        st, body, _h = await _get_text(session, front_base + src)
        if st == 200:
            chunks.append((src, body))
    return html, chunks


# Check 8: фронт отдаёт index.html и непустой JS-бандл (SLA фронта).
async def check_8(session, front_base):
    try:
        html, chunks = await _fetch_bundle(session, front_base)
        if html is None:
            return DOWN, "frontend index.html unreachable"
        if '<div id="root"' not in html or "<script" not in html:
            return CORRUPT, "frontend index.html looks broken"
        if not chunks:
            return CORRUPT, "no /assets/*.js chunks referenced"
        if max(len(body) for _, body in chunks) < 10_000:
            return CORRUPT, "bundle is unreasonably small"
        print(f"[+] [Check-8]: frontend bundle OK - {len(chunks)} chunks")
        return OK, "Check 8 - OK"
    except Exception:
        return DOWN, "frontend unreachable"


# Check 9: SPA-роуты отдают 200 (nginx try_files fallback работает).
async def check_9(session, front_base):
    try:
        for route in ("/login", "/cases", "/archive", "/hearings"):
            status, html, _ = await _get_text(session, front_base + route)
            if status != 200:
                return MUMBLE, f"{route} returned {status}"
            if '<div id="root"' not in html:
                return CORRUPT, f"{route} did not return SPA index.html"
        print("[+] [Check-9]: SPA routes OK")
        return OK, "Check 9 - OK"
    except Exception as e:
        return MUMBLE, f"SPA routes check error: {e}"


# Check 10: статические ассеты отдаются с корректным Content-Type.
async def check_10(session, front_base):
    try:
        _html, chunks = await _fetch_bundle(session, front_base)
        if not chunks:
            return CORRUPT, "no chunks to probe for headers"
        asset = chunks[0][0]
        async with session.get(front_base + asset, timeout=DEFAULT_TIMEOUT) as r:
            if r.status != 200:
                return MUMBLE, f"asset {asset} returned {r.status}"
            ct = r.headers.get("Content-Type", "")
            if "javascript" not in ct.lower():
                return CORRUPT, f"asset {asset} wrong content-type: {ct}"
        print(f"[+] [Check-10]: asset content-type OK - {ct}")
        return OK, "Check 10 - OK"
    except Exception as e:
        return MUMBLE, f"asset headers check error: {e}"


# Check 11: /api/* проксируется через nginx и отвечает 200 на /api/health.
async def check_11(session, front_base):
    try:
        status, data = await _get(session, front_base + "/api/health")
        if status != 200:
            return MUMBLE, f"/api/health via nginx returned {status}"
        if not isinstance(data, dict) or data.get("status") != "operational":
            return CORRUPT, f"/api/health via nginx payload unexpected: {data}"
        print("[+] [Check-11]: nginx -> backend proxy OK")
        return OK, "Check 11 - OK"
    except Exception:
        return DOWN, "nginx -> backend proxy unreachable"


# Check 12: e2e логин через nginx-прокси (front_base, :8081).
async def check_12(session, front_base):
    try:
        token, role = await _login(session, front_base, "citizen_07", "c1t1z3n!")
        if not token:
            return MUMBLE, "e2e login via nginx failed"
        if role != "citizen":
            return CORRUPT, f"e2e login wrong role: {role}"
        status, me = await _get(
            session, front_base + "/api/auth/me", headers=_bearer(token)
        )
        if status != 200 or not isinstance(me, dict):
            return MUMBLE, f"/api/auth/me via nginx failed: {status}"
        if me.get("username") != "citizen_07":
            return CORRUPT, f"/api/auth/me via nginx mismatch: {me}"
        print("[+] [Check-12]: e2e nginx round-trip OK")
        return OK, "Check 12 - OK"
    except Exception as e:
        return MUMBLE, f"e2e check error: {e}"


async def check(host, port):
    base = f"http://{host}:{port}"
    front_base = f"http://{host}:{FRONT_PORT}"
    verdict = OK, "OK"

    async with aiohttp.ClientSession() as session:
        checks = [
            ("Check-1", check_1, (session, base)),
            ("Check-2", check_2, (session, base)),
            ("Check-3", check_3, (session, base)),
            ("Check-4", check_4, (session, base)),
            ("Check-5", check_5, (session, base)),
            ("Check-6", check_6, (session, base)),
            ("Check-7", check_7, (session, base)),
            ("Check-8", check_8, (session, front_base)),
            ("Check-9", check_9, (session, front_base)),
            ("Check-10", check_10, (session, front_base)),
            ("Check-11", check_11, (session, front_base)),
            ("Check-12", check_12, (session, front_base)),
        ]
        for name, fn, args in checks:
            try:
                code, msg = await fn(*args)
            except Exception as e:
                code, msg = CHECKER_ERROR, f"{name}: {e}"
            if code != OK:
                return code, msg
    return verdict


async def pwn(host, port):
    try:
        verdict = await check(host, port)
        print(f"[*] Main Verdict - status: {verdict[0]}, msg: {verdict[1]}")
        return verdict
    except Exception:
        return DOWN, "Service Unavailable"


# import asyncio
# asyncio.run(pwn("127.0.0.1", 8080))
