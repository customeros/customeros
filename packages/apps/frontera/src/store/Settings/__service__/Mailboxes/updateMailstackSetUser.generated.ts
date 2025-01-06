import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type MailstackSetUserMutationVariables = Types.Exact<{
  mailbox: Types.Scalars['String']['input'];
  userId: Types.Scalars['ID']['input'];
}>;


export type MailstackSetUserMutation = { __typename?: 'Mutation', mailstack_SetUser: { __typename?: 'Result', result: boolean } };
