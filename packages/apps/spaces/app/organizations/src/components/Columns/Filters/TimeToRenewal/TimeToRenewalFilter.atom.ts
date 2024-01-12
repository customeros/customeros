import addDays from 'date-fns/addDays';
import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface TimeToRenewalState {
  value: string;
  isActive: boolean;
}

const defaultValue = addDays(new Date(), 7).toISOString().split('T')[0];

export const defaultState: TimeToRenewalState = {
  value: defaultValue,
  isActive: false,
};

export const TimeToRenewalAtom = atom<TimeToRenewalState>({
  key: 'time-to-renewal-filter',
  default: defaultState,
});

export const TimeToRenewalFilterSelector = selector({
  key: 'time-to-renewal-filter-selector',
  get: ({ get }) => get(TimeToRenewalAtom),
});

export const useTimeToRenewalFilter = () => {
  return useRecoilState(TimeToRenewalAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom TimeToRenewalState
 */
export const mapTimeToRenewalToAtom =
  makeServerToAtomMapper<TimeToRenewalState>(
    {
      filter: {
        property: 'RENEWAL_DATE',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
    }),
    defaultState,
  );
