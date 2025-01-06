import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowArchiveMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type FlowArchiveMutation = { __typename?: 'Mutation', flow_Archive: { __typename?: 'Result', result: boolean } };
