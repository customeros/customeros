import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetDomainsQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetDomainsQuery = { __typename?: 'Query', mailstack_Domains: Array<string> };
