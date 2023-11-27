import defaultsDeep from 'lodash/defaultsDeep';

import { Unpacked } from '@shared/types/helpers';
import { GetOrganizationsQuery } from '@organizations/graphql/getOrganizations.generated';
import { UpdateOrganizationMutationVariables } from '@shared/graphql/updateOrganization.generated';

export type GetOrganizationRowResult = Unpacked<
  NonNullable<GetOrganizationsQuery['dashboardView_Organizations']>['content']
>;

const defaults: GetOrganizationRowResult = {
  __typename: 'Organization',
  id: '',
  name: '',
  description: null,
  industry: null,
  website: null,
  domains: [],
  isCustomer: false,
  lastTouchPointTimelineEventId: null,
  lastTouchPointAt: null,
  subsidiaryOf: [],
  owner: {
    __typename: 'User',
    id: '',
    firstName: '',
    lastName: '',
  },
  accountDetails: {
    __typename: 'OrgAccountDetails',
    renewalForecast: {
      __typename: 'RenewalForecast',
      arr: null,
      maxArr: null,
      updatedAt: null,
      updatedById: null,
      updatedBy: {
        __typename: 'User',
        id: '',
        firstName: '',
        lastName: '',
        emails: null,
      },
    },
  },
  locations: [],
  lastTouchPointTimelineEvent: null,
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
        id: data.id,
        name: data.name,
        description: data.description,
        industry: data.industry,
        website: data.website,
        domains: data.domains,
        isCustomer: data.isCustomer,
      },
    };
  }
}
