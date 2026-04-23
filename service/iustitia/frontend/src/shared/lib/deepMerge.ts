const isObject = (v: unknown): v is Record<string, unknown> =>
  typeof v === 'object' && v !== null && !Array.isArray(v);

export const deepMerge = <T extends Record<string, unknown>>(
  target: T,
  source: Record<string, unknown>,
): T => {
  for (const key in source) {
    const srcVal = source[key];
    const tgtVal = (target as Record<string, unknown>)[key];
    if (isObject(srcVal) && isObject(tgtVal)) {
      (target as Record<string, unknown>)[key] = deepMerge(
        { ...(tgtVal as Record<string, unknown>) },
        srcVal,
      );
    } else {
      (target as Record<string, unknown>)[key] = srcVal;
    }
  }
  return target;
};
