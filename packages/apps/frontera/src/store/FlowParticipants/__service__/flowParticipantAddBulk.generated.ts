import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowParticipantAddBulkMutationVariables = Types.Exact<{
  flowId: Types.Scalars['ID']['input'];
  entitiIds: Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input'];
  entityType: Types.FlowEntityType;
}>;


export type FlowParticipantAddBulkMutation = { __typename?: 'Mutation', flowParticipant_AddBulk: { __typename?: 'Result', result: boolean } };
