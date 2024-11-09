import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FindWorkContactEmailMutationVariables = Types.Exact<{
  contactId: Types.Scalars['ID']['input'];
  organizationId: Types.Scalars['ID']['input'];
}>;

export type FindWorkContactEmailMutation = {
  __typename?: 'Mutation';
  contact_FindWorkEmail: { __typename?: 'ActionResponse'; accepted: boolean };
};
