import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateContactBulkByLinkedInMutationVariables = Types.Exact<{
  linkedInUrls: Array<Types.Scalars['String']['input']> | Types.Scalars['String']['input'];
  flowId?: Types.InputMaybe<Types.Scalars['String']['input']>;
}>;


export type CreateContactBulkByLinkedInMutation = { __typename?: 'Mutation', contact_CreateBulkByLinkedIn: Array<string> };
