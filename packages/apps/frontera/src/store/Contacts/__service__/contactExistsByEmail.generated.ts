import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ContactByEmailQueryVariables = Types.Exact<{
  email: Types.Scalars['String']['input'];
}>;

export type ContactByEmailQuery = {
  __typename?: 'Query';
  contact_ByEmail?: {
    __typename?: 'Contact';
    emails: Array<{ __typename?: 'Email'; email?: string | null }>;
  } | null;
};
