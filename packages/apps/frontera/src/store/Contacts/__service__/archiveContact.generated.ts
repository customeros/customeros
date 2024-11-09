import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ArchiveContactMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
}>;

export type ArchiveContactMutation = {
  __typename?: 'Mutation';
  contact_Hide: { __typename?: 'ActionResponse'; accepted: boolean };
};
