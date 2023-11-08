import { selector, useRecoilValue } from 'recoil';

import { ForecastFilterSelector } from '@organizations/components/Columns/Filters/Forecast/ForecastFilter.atom';
import { OrganizationFilterSelector } from '@organizations/components/Columns/Filters/Organization/OrganizationFilter.atom';
import { RelationshipFilterSelector } from '@organizations/components/Columns/Filters/Relationship/RelationshipFilter.atom';
import { TimeToRenewalFilterSelector } from '@organizations/components/Columns/Filters/TimeToRenewal/TimeToRenewalFilter.atom';
import { RenewalLikelihoodFilterSelector } from '@organizations/components/Columns/Filters/RenewalLikelihood/RenewalLikelihoodFilter.atom';

const tableStateSelector = selector({
  key: 'tableState',
  get: ({ get }) => {
    const forecast = get(ForecastFilterSelector);
    const organization = get(OrganizationFilterSelector);
    const relationship = get(RelationshipFilterSelector);
    const timeToRenewal = get(TimeToRenewalFilterSelector);
    const renewalLikelihood = get(RenewalLikelihoodFilterSelector);

    return {
      columnFilters: {
        forecast,
        organization,
        relationship,
        timeToRenewal,
        renewalLikelihood,
      },
    };
  },
});

export const useTableState = () => {
  return useRecoilValue(tableStateSelector);
};
