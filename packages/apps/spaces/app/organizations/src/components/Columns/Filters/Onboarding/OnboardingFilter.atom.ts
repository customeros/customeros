import { atom, selector, useRecoilState } from 'recoil';

import { OnboardingStatus } from '@graphql/types';

interface OnboardingFilterState {
  isActive: boolean;
  value: OnboardingStatus[];
}

export const defaultState: OnboardingFilterState = {
  value: [],
  isActive: false,
};

export const OnboardingFilterAtom = atom<OnboardingFilterState>({
  key: 'onboarding-filter',
  default: defaultState,
});

export const OnboardingFilterSelector = selector({
  key: 'onboarding-filter-selector',
  get: ({ get }) => get(OnboardingFilterAtom),
});

export const useOnboardingFilter = () => {
  return useRecoilState(OnboardingFilterAtom);
};
