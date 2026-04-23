// Patch F4 (Prototype Pollution).
// Raw-версия реализовывала рекурсивный deepMerge руками и копировала
// ключи '__proto__', 'prototype', 'constructor' в target без фильтрации.
// Вход - query-параметр ?config= в entities/archive/model/filters.ts
// (JSON.parse пользовательского payload -> deepMerge), т.е. злоумышленник
// может отравить Object.prototype через
//   /archive?config={"__proto__":{"isAdmin":true}}
// и получить глобальное isAdmin: true на всех объектах в рантайме.
//
// Фикс состоит из двух слоёв:
//   1. sanitizeKeys() - рекурсивно выкидываем FORBIDDEN_KEYS до merge-а.
//   2. lodash.merge (well-maintained, де-факто стандарт) - его внутренний
//      обход тоже игнорирует "__proto__"/"prototype" как пути (см. CVE
//      CVE-2019-10744, закрытый в lodash 4.17.12), что даёт defence-in-depth.
// Экспортная сигнатура сохранена, чтобы не трогать вызывающий код в
// filters.ts.
import merge from 'lodash.merge';

const FORBIDDEN_KEYS = new Set(['__proto__', 'prototype', 'constructor']);

const isPlainObject = (v: unknown): v is Record<string, unknown> =>
  typeof v === 'object' && v !== null && !Array.isArray(v);

const sanitizeKeys = (value: unknown): unknown => {
  if (Array.isArray(value)) {
    return value.map(sanitizeKeys);
  }
  if (!isPlainObject(value)) {
    return value;
  }
  const out: Record<string, unknown> = {};
  for (const [key, val] of Object.entries(value)) {
    if (FORBIDDEN_KEYS.has(key)) continue;
    out[key] = sanitizeKeys(val);
  }
  return out;
};

export const deepMerge = <T extends Record<string, unknown>>(
  target: T,
  source: Record<string, unknown>,
): T => {
  const safeSource = sanitizeKeys(source) as Record<string, unknown>;
  return merge({}, target, safeSource) as T;
};
