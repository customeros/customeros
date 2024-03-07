import { selector, useRecoilValue } from 'recoil';
import { OwnerFilterSelector } from '@renewals/components/Columns/Filters/Owner/OwnerFilter.atom';
import { ForecastFilterSelector } from '@renewals/components/Columns/Filters/Forecast/ForecastFilter.atom';
import { OrganizationFilterSelector } from '@renewals/components/Columns/Filters/Organization/OrganizationFilter.atom';
import { LastTouchpointSelector } from '@renewals/components/Columns/Filters/LastTouchpoint/LastTouchpointFilter.atom';
import { TimeToRenewalFilterSelector } from '@renewals/components/Columns/Filters/TimeToRenewal/TimeToRenewalFilter.atom';
import { RenewalLikelihoodFilterSelector } from '@renewals/components/Columns/Filters/RenewalLikelihood/RenewalLikelihoodFilter.atom';

const tableStateSelector = selector({
  key: 'renewalsTableState',
  get: ({ get }) => {
    const owner = get(OwnerFilterSelector);
    const forecast = get(ForecastFilterSelector);
    const organization = get(OrganizationFilterSelector);
    const renewalLikelihood = get(RenewalLikelihoodFilterSelector);

    const timeToRenewal = (() => {
      const state = get(TimeToRenewalFilterSelector);
      const value = new Date(state.value).toISOString();

      return {
        ...state,
        value,
      };
    })();

    const lastTouchpoint = (() => {
      const state = get(LastTouchpointSelector);
      const after = state.after
        ? new Date(state.after).toISOString()
        : undefined;

      return {
        ...state,
        after,
      };
    })();

    return {
      columnFilters: {
        owner,
        forecast,
        organization,
        timeToRenewal,
        lastTouchpoint,
        renewalLikelihood,
      },
    };
  },
});

export const useTableState = () => {
  return useRecoilValue(tableStateSelector);
};
