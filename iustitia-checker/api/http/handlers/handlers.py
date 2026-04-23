import time

import requests
import uvicorn
from fastapi import APIRouter, FastAPI

# Импорт основных скриптов для проверок сервиса
from scripts.iustitia.checkers.checker import pwn as checker
from scripts.iustitia.sploits.sploit import pwn as sploit


class Server:
    def __init__(
        self,
        host: str,
        port: int,
        operations_hostport: str,
        standalone: bool = False,
    ):
        self.operations_hostport = operations_hostport
        self.host = host
        self.port = port
        self.standalone = standalone

        self.app = FastAPI()
        self.router = self._create_router()
        self.app.include_router(self.router)

        self.health_ok = False
        self.last_response = None

    def _register(self):
        if self.standalone or not self.operations_hostport:
            print("[checker] standalone mode - skipping operations registration")
            self.health_ok = True
            return

        url = "http://{}/game/checker/register".format(self.operations_hostport)
        timeout = 5.0
        delay = 10

        while not self.health_ok:
            params = {"vuln_service": "iustitia", "serve_port": self.port}
            response = requests.post(url, params=params, timeout=timeout)
            self.last_response = {
                "status_code": response.status_code,
            }

            try:
                response.raise_for_status()

                if response.status_code == 200:
                    self.health_ok = True
                    print("External API health check PASSED (200)")
                else:
                    print(
                        "External API returned %d, retrying in %d sec...",
                        response.status_code,
                        delay,
                    )

            except Exception as e:
                print(f"Health check failed: {str(e)}")
            finally:
                response.close()

            if not self.health_ok:
                time.sleep(delay)

    def _create_router(self) -> APIRouter:
        router = APIRouter(prefix="/api", tags=["Api"])

        @router.post("/check")
        async def start_check(vulnbox_host: str, vulnbox_port: int):
            check = await checker(vulnbox_host, vulnbox_port)

            if check[0] == 101:
                print("[+] [Checker]: activate sploits...")
                vuln = await sploit(vulnbox_host, vulnbox_port)

                return {
                    "data": {
                        "status_code": check[0],
                        "check_message": check[1],
                        "fixed_vulns": vuln,
                    }
                }
            else:
                return {
                    "data": {
                        "status_code": check[0],
                        "check_message": check[1],
                        "fixed_vulns": 0,
                    }
                }

        return router

    def run(self):
        self._register()
        uvicorn.run(self.app, host=self.host, port=self.port)
