import type { Case } from '@shared/api';
import { create } from 'zustand';

export type StatusFilter = 'all' | Case['status'];

interface CaseFilterState {
  status: StatusFilter;
  query: string;
  setStatus: (status: StatusFilter) => void;
  setQuery: (query: string) => void;
}

export const useCaseFilterStore = create<CaseFilterState>((set) => ({
  status: 'all',
  query: '',
  setStatus: (status) => set({ status }),
  setQuery: (query) => set({ query }),
}));
