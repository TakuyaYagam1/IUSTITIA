declare const __INTERNAL_HMAC_KEY__: string;

export const INTERNAL_HMAC_KEY: string =
  typeof __INTERNAL_HMAC_KEY__ === 'string' ? __INTERNAL_HMAC_KEY__ : '';

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
