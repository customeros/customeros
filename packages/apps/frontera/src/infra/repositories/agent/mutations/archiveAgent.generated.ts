import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type AgentArchiveMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type AgentArchiveMutation = {
  __typename?: 'Mutation';
  agent_Delete: boolean;
};
