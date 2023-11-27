import { selector, useRecoilValue } from 'recoil';

import { OwnerFilterSelector } from '@organizations/components/Columns/Filters/Owner/OwnerFilter.atom';
import { WebsiteFilterSelector } from '@organizations/components/Columns/Filters/Website/WebsiteFilter.atom';
import { ForecastFilterSelector } from '@organizations/components/Columns/Filters/Forecast/ForecastFilter.atom';
import { OrganizationFilterSelector } from '@organizations/components/Columns/Filters/Organization/OrganizationFilter.atom';
import { RelationshipFilterSelector } from '@organizations/components/Columns/Filters/Relationship/RelationshipFilter.atom';
import { LastTouchpointSelector } from '@organizations/components/Columns/Filters/LastTouchpoint/LastTouchpointFilter.atom';

const tableStateSelector = selector({
  key: 'tableState',
  get: ({ get }) => {
    const owner = get(OwnerFilterSelector);
    const website = get(WebsiteFilterSelector);
    const forecast = get(ForecastFilterSelector);
    const organization = get(OrganizationFilterSelector);
    const relationship = get(RelationshipFilterSelector);

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
        website,
        forecast,
        organization,
        relationship,
        lastTouchpoint,
      },
    };
  },
});

export const useTableState = () => {
  return useRecoilValue(tableStateSelector);
};
