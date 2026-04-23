import { deepMerge } from '@shared/lib/deepMerge';

export interface ClassifiedConfig {
  verdict: 'all' | 'guilty' | 'acquitted';
  sortBy: 'archived_at' | 'defendant';
  order: 'asc' | 'desc';
  limit: number;
  [key: string]: unknown;
}

export const DEFAULT_CONFIG: ClassifiedConfig = {
  verdict: 'all',
  sortBy: 'archived_at',
  order: 'desc',
  limit: 50,
};

const STORAGE_KEY = 'archive_config';

export const loadClassifiedConfig = (search: string): ClassifiedConfig => {
  const base: ClassifiedConfig = { ...DEFAULT_CONFIG };

  let fromStorage: Record<string, unknown> = {};
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) fromStorage = JSON.parse(raw) as Record<string, unknown>;
  } catch {
    fromStorage = {};
  }

  const params = new URLSearchParams(search);
  let fromQuery: Record<string, unknown> = {};
  const cfgParam = params.get('config');
  if (cfgParam) {
    try {
      fromQuery = JSON.parse(cfgParam) as Record<string, unknown>;
    } catch {
      fromQuery = {};
    }
  }

  const merged = deepMerge(base, fromStorage);
  return deepMerge(merged, fromQuery);
};

export const saveClassifiedConfig = (cfg: ClassifiedConfig): void => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(cfg));
  } catch {
    /* noop */
  }
};
