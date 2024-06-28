import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateReminderMutationVariables = Types.Exact<{
  input: Types.ReminderUpdateInput;
}>;


export type UpdateReminderMutation = { __typename?: 'Mutation', reminder_Update?: string | null };
