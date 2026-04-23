// Patch F5 (Secrets leak в JS-бандл, HMAC-ключ).
// Raw-версия делала то же, что service.ts: объявляла
// __INTERNAL_HMAC_KEY__ через vite define: и читала из него
// SECRET_MARKER_K_*. Значение использовалось для подписания выдуманных
// "/api/internal/*" запросов; бек такого роута не имеет, поэтому фронту
// ключ не нужен вообще. Фикс - пустой литерал и удаление define:
// и build-args.
export const INTERNAL_HMAC_KEY: string = '';

const hexEncode = (buf: ArrayBuffer): string => {
  const bytes = new Uint8Array(buf);
  let hex = '';
  for (const b of bytes) hex += b.toString(16).padStart(2, '0');
  return hex;
};

export const signInternalRequest = async (payload: string): Promise<string> => {
  const keyBytes = new TextEncoder().encode(INTERNAL_HMAC_KEY);
  const msgBytes = new TextEncoder().encode(payload);
  const key = await crypto.subtle.importKey(
    'raw',
    keyBytes,
    { name: 'HMAC', hash: 'SHA-256' },
    false,
    ['sign'],
  );
  const sig = await crypto.subtle.sign('HMAC', key, msgBytes);
  return hexEncode(sig);
};
