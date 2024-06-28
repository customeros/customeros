import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type CreateLogEntryMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  logEntry: Types.LogEntryInput;
}>;

export type CreateLogEntryMutation = {
  __typename?: 'Mutation';
  logEntry_CreateForOrganization: string;
};
