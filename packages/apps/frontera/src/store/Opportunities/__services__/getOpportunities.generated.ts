import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetOpportunitiesQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
}>;

export type GetOpportunitiesQuery = {
  __typename?: 'Query';
  opportunities_LinkedToOrganizations: {
    __typename?: 'OpportunityPage';
    totalElements: any;
    totalAvailable: any;
    content: Array<{
      __typename?: 'Opportunity';
      name: string;
      amount: number;
      maxAmount: number;
      internalType: Types.InternalType;
      externalType: string;
      internalStage: Types.InternalStage;
      externalStage: string;
      estimatedClosedAt?: any | null;
      generalNotes: string;
      nextSteps: string;
      renewedAt?: any | null;
      currency?: Types.Currency | null;
      stageLastUpdated?: any | null;
      renewalApproved: boolean;
      renewalLikelihood: Types.OpportunityRenewalLikelihood;
      renewalUpdatedByUserId: string;
      renewalUpdatedByUserAt?: any | null;
      renewalAdjustedRate: any;
      comments: string;
      likelihoodRate: any;
      id: string;
      createdAt?: any | null;
      updatedAt?: any | null;
      source?: Types.DataSource | null;
      appSource?: string | null;
      metadata: {
        __typename?: 'Metadata';
        id: string;
        created: any;
        lastUpdated: any;
        source: Types.DataSource;
        sourceOfTruth: Types.DataSource;
        appSource: string;
      };
      organization?: {
        __typename?: 'Organization';
        metadata: {
          __typename?: 'Metadata';
          id: string;
          created: any;
          lastUpdated: any;
          sourceOfTruth: Types.DataSource;
        };
      } | null;
      createdBy?: {
        __typename?: 'User';
        id: string;
        firstName: string;
        lastName: string;
        name?: string | null;
      } | null;
      owner?: {
        __typename?: 'User';
        id: string;
        firstName: string;
        lastName: string;
        name?: string | null;
        appSource: string;
        bot: boolean;
        createdAt: any;
        calendars: Array<{ __typename?: 'Calendar'; id: string }>;
      } | null;
      externalLinks: Array<{
        __typename?: 'ExternalSystem';
        externalUrl?: string | null;
        externalId?: string | null;
      }>;
    }>;
  };
};
