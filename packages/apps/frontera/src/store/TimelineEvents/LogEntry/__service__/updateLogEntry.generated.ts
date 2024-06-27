import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type UpdateLogEntryMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  input: Types.LogEntryUpdateInput;
}>;

export type UpdateLogEntryMutation = {
  __typename?: 'Mutation';
  logEntry_Update: string;
};
