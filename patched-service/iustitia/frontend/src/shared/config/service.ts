// Patch F5 (Secrets leak в JS-бандл).
// Raw-версия:
//   declare const __SERVICE_TOKEN__: string;
//   export const SERVICE_TOKEN = __SERVICE_TOKEN__;
// Литерал __SERVICE_TOKEN__ инлайнился Vite через define: {} (см.
// vite.config.ts) из VITE_SERVICE_TOKEN, который шёл через docker-compose
// build-args и Dockerfile ARG/ENV. В итоге "SECRET_MARKER_B_..." попадал
// как string-literal в /assets/index-*.js и раздавался nginx-ом любому
// curl-клиенту. Фикс - полностью убираем define: и читаем из пустоты
// (фронт реально не обращается к этому эндпоинту, см. shared/api/client.ts,
// где токен пробрасывается только в '/api/internal/*', а такого роута на
// беке нет).
export const SERVICE_TOKEN: string = '';
