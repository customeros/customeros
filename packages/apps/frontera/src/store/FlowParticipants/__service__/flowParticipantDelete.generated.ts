import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowParticipantDeleteMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type FlowParticipantDeleteMutation = { __typename?: 'Mutation', flowParticipant_Delete: { __typename?: 'Result', result: boolean } };
