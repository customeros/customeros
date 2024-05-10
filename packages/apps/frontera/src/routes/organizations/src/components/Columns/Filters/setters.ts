import { useEffect } from 'react';

import { match } from 'ts-pattern';
import { useRecoilTransaction_UNSTABLE } from 'recoil';

import { Filter } from '@graphql/types';

import { getServerToAtomMapper } from './serverMappers';
import { OwnerFilterAtom, OwnerFilterState } from './Owner';
import { WebsiteFilterAtom, WebsiteFilterState } from './Website';
import { ForecastFilterAtom, ForecastFilterState } from './Forecast';
import { TimeToRenewalAtom, TimeToRenewalState } from './TimeToRenewal';
import { OnboardingFilterAtom, OnboardingFilterState } from './Onboarding';
import { LastTouchpointAtom, LastTouchpointState } from './LastTouchpoint';
import {
  OrganizationFilterAtom,
  OrganizationFilterState,
} from './Organization';
import {
  RelationshipFilterAtom,
  RelationshipFilterState,
} from './Relationship';
import {
  RenewalLikelihoodFilterAtom,
  RenewalLiklihoodFilterState,
} from './RenewalLikelihood';

export const parseRawFilters = (raw = '') => {
  if (!raw) return [];
  const filterData = JSON.parse(raw) as { filter: Filter };

  if (!filterData?.filter?.AND) return [];

  return filterData?.filter?.AND?.map((data) => {
    const property = data?.filter?.property ?? '';
    const mapToAtom = getServerToAtomMapper(property);

    if (mapToAtom) {
      return [property, mapToAtom(data)];
    }
  }).filter(Boolean) as [string, ForecastFilterState | LastTouchpointState][];
};

export const useFilterSetter = (rawFilters?: string | null) => {
  const parsedFilters = parseRawFilters(rawFilters ?? '');

  const setFilters = useRecoilTransaction_UNSTABLE(
    ({ set }) =>
      (id: string, value: unknown) => {
        match(id)
          .with('FORECAST_ARR', () => {
            set<ForecastFilterState>(
              ForecastFilterAtom,
              value as ForecastFilterState,
            );
          })
          .with('LAST_TOUCHPOINT_TYPE', 'LAST_TOUCHPOINT_AT', () => {
            set<LastTouchpointState>(
              LastTouchpointAtom,
              value as LastTouchpointState,
            );
          })
          .with('OWNER_ID', () => {
            set<OwnerFilterState>(OwnerFilterAtom, value as OwnerFilterState);
          })
          .with('WEBSITE', () => {
            set<WebsiteFilterState>(
              WebsiteFilterAtom,
              value as WebsiteFilterState,
            );
          })
          .with('ONBOARDING_STATUS', () => {
            set<OnboardingFilterState>(
              OnboardingFilterAtom,
              value as OnboardingFilterState,
            );
          })
          .with('RENEWAl_DATE', () => {
            set<TimeToRenewalState>(
              TimeToRenewalAtom,
              value as TimeToRenewalState,
            );
          })
          .with('IS_CUSTOMER', () => {
            set<RelationshipFilterState>(
              RelationshipFilterAtom,
              value as RelationshipFilterState,
            );
          })
          .with('NAME', () => {
            set<OrganizationFilterState>(
              OrganizationFilterAtom,
              value as OrganizationFilterState,
            );
          })
          .with('RENEWAL_LIKELIHOOD', () => {
            set<RenewalLiklihoodFilterState>(
              RenewalLikelihoodFilterAtom,
              value as RenewalLiklihoodFilterState,
            );
          })
          .otherwise(() => {});
      },
  );

  useEffect(() => {
    parsedFilters.forEach(([id, value]) => {
      setFilters(id, value);
    });
  }, [rawFilters]);
};
