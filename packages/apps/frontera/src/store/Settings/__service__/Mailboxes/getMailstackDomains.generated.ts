import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type MailstackDomainsQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type MailstackDomainsQuery = {
  __typename?: 'Query';
  mailstack_Domains: Array<string>;
};
