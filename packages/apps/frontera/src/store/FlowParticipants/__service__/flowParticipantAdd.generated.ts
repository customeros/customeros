import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowParticipantAddMutationVariables = Types.Exact<{
  flowId: Types.Scalars['ID']['input'];
  entityId: Types.Scalars['ID']['input'];
  entityType: Types.FlowEntityType;
}>;


export type FlowParticipantAddMutation = { __typename?: 'Mutation', flowParticipant_Add: { __typename?: 'FlowParticipant', metadata: { __typename?: 'Metadata', id: string } } };
