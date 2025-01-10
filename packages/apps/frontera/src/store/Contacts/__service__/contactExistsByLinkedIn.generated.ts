import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ContactExistsByLinkedInQueryVariables = Types.Exact<{
  linkedIn: Types.Scalars['String']['input'];
}>;

export type ContactExistsByLinkedInQuery = {
  __typename?: 'Query';
  contact_ByLinkedIn?: {
    __typename?: 'Contact';
    metadata: { __typename?: 'Metadata'; id: string };
  } | null;
};
