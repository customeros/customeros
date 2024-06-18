import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type CreateReminderMutationVariables = Types.Exact<{
  input: Types.ReminderInput;
}>;

export type CreateReminderMutation = {
  __typename?: 'Mutation';
  reminder_Create?: string | null;
};
