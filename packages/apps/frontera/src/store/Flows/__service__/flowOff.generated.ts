import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowOffMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type FlowOffMutation = { __typename?: 'Mutation', flow_Off: { __typename?: 'Flow', metadata: { __typename?: 'Metadata', id: string } } };
