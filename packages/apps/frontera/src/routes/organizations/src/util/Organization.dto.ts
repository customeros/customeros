import defaultsDeep from 'lodash/defaultsDeep';

import { Unpacked } from '@shared/types/helpers';
import { GetOrganizationsQuery } from '@organizations/graphql/getOrganizations.generated';
import { UpdateOrganizationMutationVariables } from '@shared/graphql/updateOrganization.generated';

export type GetOrganizationRowResult = Unpacked<
  NonNullable<GetOrganizationsQuery['dashboardView_Organizations']>['content']
>;

const defaults: GetOrganizationRowResult = {
  metadata: {
    id: '',
    created: new Date().toISOString(),
  },
  name: '',
  description: null,
  industry: null,
  website: null,
  domains: [],
  relationship: undefined,
  lastTouchpoint: {
    lastTouchPointTimelineEventId: null,
    lastTouchPointAt: null,
  },
  parentCompanies: [],
  owner: {
    __typename: 'User',
    id: '',
    firstName: '',
    lastName: '',
  },
  socialMedia: [],
  accountDetails: {
    __typename: 'OrgAccountDetails',

    renewalSummary: {
      __typename: 'RenewalSummary',
      arrForecast: null,
      maxArrForecast: null,
      renewalLikelihood: null,
      nextRenewalDate: null,
    },
  },
  locations: [],
};

export class OrganizationRowDTO {
  constructor(data: GetOrganizationRowResult) {
    Object.assign(this, defaultsDeep(data, defaults));
  }

  static toUpdatePayload(
    data: GetOrganizationRowResult,
  ): UpdateOrganizationMutationVariables {
    return {
      input: {
        id: data.metadata.id,
        name: data.name,
        description: data.description,
        industry: data.industry,
        website: data.website,
        domains: data.domains,
        relationship: data.relationship,
      },
    };
  }
}
