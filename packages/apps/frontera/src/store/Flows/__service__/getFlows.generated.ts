import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetFlowsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GetFlowsQuery = {
  __typename?: 'Query';
  flows: Array<{
    __typename?: 'Flow';
    name: string;
    description: string;
    status: Types.FlowStatus;
    metadata: {
      __typename?: 'Metadata';
      id: string;
      created: any;
      lastUpdated: any;
      source: Types.DataSource;
      sourceOfTruth: Types.DataSource;
      appSource: string;
    };
    sequences: Array<{
      __typename?: 'FlowSequence';
      metadata: { __typename?: 'Metadata'; id: string };
    }>;
  }>;
};
