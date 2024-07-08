import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type AddTagToLogEntryMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  input: Types.TagIdOrNameInput;
}>;

export type AddTagToLogEntryMutation = {
  __typename?: 'Mutation';
  logEntry_AddTag: string;
};
