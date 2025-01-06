import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateContactBulkByEmailMutationVariables = Types.Exact<{
  emails: Array<Types.Scalars['String']['input']> | Types.Scalars['String']['input'];
  flowId?: Types.InputMaybe<Types.Scalars['String']['input']>;
}>;


export type CreateContactBulkByEmailMutation = { __typename?: 'Mutation', contact_CreateBulkByEmail: Array<string> };
