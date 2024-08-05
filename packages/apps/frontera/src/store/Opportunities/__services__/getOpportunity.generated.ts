import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type OpportunityQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type OpportunityQuery = {
  __typename?: 'Query';
  opportunity?: {
    __typename?: 'Opportunity';
    id: string;
    createdAt?: any | null;
    updatedAt?: any | null;
    name: string;
    amount: number;
    currency?: Types.Currency | null;
    maxAmount: number;
    internalType: Types.InternalType;
    externalType: string;
    internalStage: Types.InternalStage;
    externalStage: string;
    stageLastUpdated?: any | null;
    estimatedClosedAt?: any | null;
    generalNotes: string;
    nextSteps: string;
    renewedAt?: any | null;
    renewalApproved: boolean;
    renewalLikelihood: Types.OpportunityRenewalLikelihood;
    renewalUpdatedByUserId: string;
    renewalUpdatedByUserAt?: any | null;
    renewalAdjustedRate: any;
    comments: string;
    source?: Types.DataSource | null;
    sourceOfTruth?: Types.DataSource | null;
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
    } | null;
    owner?: {
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
    } | null;
    externalLinks: Array<{
      __typename?: 'ExternalSystem';
      externalId?: string | null;
      externalUrl?: string | null;
    }>;
  } | null;
};
