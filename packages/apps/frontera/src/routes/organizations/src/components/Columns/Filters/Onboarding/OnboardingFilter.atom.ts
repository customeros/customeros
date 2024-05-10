import { atom, selector, useRecoilState } from 'recoil';

import { OnboardingStatus } from '@graphql/types';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface OnboardingFilterState {
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

/**
 * Used for mapping server-side Filter data to client-side atom OnboardingFilterState
 */
export const mapOnboardingToAtom =
  makeServerToAtomMapper<OnboardingFilterState>(
    {
      filter: {
        property: 'ONBOARDING_STATUS',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
    }),
    defaultState,
  );
