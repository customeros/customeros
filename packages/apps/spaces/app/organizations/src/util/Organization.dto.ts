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
      amount: null,
      potentialAmount: null,
      comment: null,
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
    renewalLikelihood: {
      __typename: 'RenewalLikelihood',
      probability: null,
      previousProbability: null,
      comment: null,
      updatedById: null,
      updatedAt: null,
      updatedBy: {
        __typename: 'User',
        id: '',
        firstName: '',
        lastName: '',
        emails: null,
      },
    },
    billingDetails: {
      __typename: 'BillingDetails',
      renewalCycle: null,
      frequency: null,
      amount: null,
      renewalCycleNext: null,
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
