import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowChangeNameMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  name: Types.Scalars['String']['input'];
}>;


export type FlowChangeNameMutation = { __typename?: 'Mutation', flow_ChangeName: { __typename?: 'Flow', name: string, metadata: { __typename?: 'Metadata', id: string } } };
