import { atom, selector, useRecoilState } from 'recoil';

interface WebsiteFilterState {
  value: string;
  isActive: boolean;
}

export const WebsiteFilterAtom = atom<WebsiteFilterState>({
  key: 'website-filter',
  default: {
    value: '',
    isActive: false,
  },
});

export const WebsiteFilterSelector = selector({
  key: 'website-filter-selector',
  get: ({ get }) => get(WebsiteFilterAtom),
});

export const useWebsiteFilter = () => {
  return useRecoilState(WebsiteFilterAtom);
};
