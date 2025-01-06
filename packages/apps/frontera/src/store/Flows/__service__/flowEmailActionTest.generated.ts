import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowEmailActionTestMutationVariables = Types.Exact<{
  subject: Types.Scalars['String']['input'];
  bodyTemplate: Types.Scalars['String']['input'];
  sendToEmailAddress: Types.Scalars['String']['input'];
}>;


export type FlowEmailActionTestMutation = { __typename?: 'Mutation', flowEmailActionTest: { __typename?: 'Result', result: boolean } };
