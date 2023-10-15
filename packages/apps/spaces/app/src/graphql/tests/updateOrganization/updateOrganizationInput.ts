import {
  FundingRound,
  Market,
  OrganizationUpdateInput,
} from '../../../../../app/src/types/__generated__/graphql.types';

export const updateOrganizationInputVariables: OrganizationUpdateInput = {
  id: '',
  patch: true,
  name: 'updateOrganizationName',
  description: 'updateOrganizationDescription',
  note: 'updateOrganizationNote',
  domains: [],
  website: 'www.updateOrganization.com',
  industry: 'updateOrganizationIndustry',
  subIndustry: 'updateOrganizationSubIndustry',
  industryGroup: 'updateOrganizationIndustryGroup',
  isPublic: true,
  isCustomer: true,
  market: Market.B2B,
  employees: 3,
  targetAudience: 'updateOrganizationTargetAudience',
  valueProposition: 'updateOrganizationValueProposition',
  lastFundingRound: FundingRound.SeriesA,
  lastFundingAmount: 'updateOrganizationLastFundingAmount',
};
