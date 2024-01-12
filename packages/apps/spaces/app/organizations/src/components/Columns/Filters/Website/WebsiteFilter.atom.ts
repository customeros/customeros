import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface WebsiteFilterState {
  value: string;
  isActive: boolean;
  showEmpty: boolean;
}

export const defaultState: WebsiteFilterState = {
  value: '',
  isActive: false,
  showEmpty: false,
};

export const WebsiteFilterAtom = atom<WebsiteFilterState>({
  key: 'website-filter',
  default: defaultState,
});

export const WebsiteFilterSelector = selector({
  key: 'website-filter-selector',
  get: ({ get }) => get(WebsiteFilterAtom),
});

export const useWebsiteFilter = () => {
  return useRecoilState(WebsiteFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom WebsiteFilterState
 */
export const mapWebsiteToAtom = makeServerToAtomMapper<WebsiteFilterState>(
  {
    filter: {
      property: 'WEBSITE',
    },
  },
  ({ filter }) => ({
    isActive: true,
    value: filter?.value,
    showEmpty: filter?.includeEmpty ?? false,
  }),
  defaultState,
);
