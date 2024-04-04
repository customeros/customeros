import addDays from 'date-fns/addDays';
import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '@shared/components/Filters/makeServerToAtomMapper';

export interface IssueDateFilterState {
  value: string;
  isActive: boolean;
}

const defaultValue = addDays(new Date(), 7).toISOString().split('T')[0];

export const defaultState: IssueDateFilterState = {
  value: defaultValue,
  isActive: false,
};

export const IssueDateFilterAtom = atom<IssueDateFilterState>({
  key: 'issue-date-filter',
  default: defaultState,
});

export const IssueDateFilterSelector = selector({
  key: 'issue-date-filter-selector',
  get: ({ get }) => get(IssueDateFilterAtom),
});

export const useIssueDateFilter = () => {
  return useRecoilState(IssueDateFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom IssueDateState
 */
export const mapIssueDateToAtom = makeServerToAtomMapper<IssueDateFilterState>(
  {
    filter: {
      property: 'INVOICE_ISSUE_DATE',
    },
  },
  ({ filter }) => ({
    isActive: true,
    value: filter?.value,
  }),
  defaultState,
);
