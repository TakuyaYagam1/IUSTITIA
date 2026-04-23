declare const __SERVICE_TOKEN__: string;

export const SERVICE_TOKEN: string = typeof __SERVICE_TOKEN__ === 'string' ? __SERVICE_TOKEN__ : '';
