import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowOnMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type FlowOnMutation = { __typename?: 'Mutation', flow_On: { __typename?: 'Flow', metadata: { __typename?: 'Metadata', id: string } } };
