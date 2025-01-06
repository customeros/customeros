import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowParticipantDeleteBulkMutationVariables = Types.Exact<{
  id: Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input'];
}>;


export type FlowParticipantDeleteBulkMutation = { __typename?: 'Mutation', flowParticipant_DeleteBulk: { __typename?: 'Result', result: boolean } };
