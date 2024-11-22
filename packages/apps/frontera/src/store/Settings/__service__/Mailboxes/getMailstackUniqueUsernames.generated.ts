import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type MailstackUniqueUsernamesQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type MailstackUniqueUsernamesQuery = {
  __typename?: 'Query';
  mailstack_UniqueUsernames: Array<string>;
};
