import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type RemoveTagFromLogEntryMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  input: Types.TagIdOrNameInput;
}>;

export type RemoveTagFromLogEntryMutation = {
  __typename?: 'Mutation';
  logEntry_RemoveTag: string;
};
