import { P } from 'ts-pattern';
import subDays from 'date-fns/subDays';
import { atom, selector, useRecoilState, useSetRecoilState } from 'recoil';

import { touchpoints } from './util';
import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface LastTouchpointState {
  after?: string;
  value: string[];
  isActive: boolean;
}

const defaultAfter = subDays(new Date(), 7).toISOString().split('T')[0];

export const defaultState: LastTouchpointState = {
  value: touchpoints.map((t) => t.value),
  after: defaultAfter,
  isActive: false,
};

export const LastTouchpointAtom = atom<LastTouchpointState>({
  key: 'last-touchpoint-filter',
  default: defaultState,
});

export const LastTouchpointSelector = selector({
  key: 'last-touchpoint-filter-selector',
  get: ({ get }) => get(LastTouchpointAtom),
});

export const useLastTouchpointFilter = () => {
  return useRecoilState(LastTouchpointAtom);
};

export const useSetLastTouchpointFilter = () => {
  return useSetRecoilState(LastTouchpointAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom LastTouchpointState
 */
export const mapLastTouchpointToAtom =
  makeServerToAtomMapper<LastTouchpointState>(
    {
      filter: {
        property: P.union('LAST_TOUCHPOINT_TYPE', 'LAST_TOUCHPOINT_AT'),
      },
    },
    ({ filter }) =>
      filter?.property === 'LAST_TOUCHPOINT_TYPE'
        ? {
            isActive: true,
            value: filter?.value,
            after: defaultAfter,
          }
        : {
            isActive: true,
            value: touchpoints.map((t) => t.value),
            after: filter?.value,
          },
    defaultState,
  );
