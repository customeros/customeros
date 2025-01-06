import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowArchiveBulkMutationVariables = Types.Exact<{
  ids: Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input'];
}>;


export type FlowArchiveBulkMutation = { __typename?: 'Mutation', flow_ArchiveBulk: { __typename?: 'Result', result: boolean } };
