import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowSequenceUnlinkContactMutationVariables = Types.Exact<{
  sequenceId: Types.Scalars['ID']['input'];
  contactId: Types.Scalars['ID']['input'];
  emailId: Types.Scalars['ID']['input'];
}>;

export type FlowSequenceUnlinkContactMutation = {
  __typename?: 'Mutation';
  flow_sequence_UnlinkContact: { __typename?: 'Result'; result: boolean };
};
