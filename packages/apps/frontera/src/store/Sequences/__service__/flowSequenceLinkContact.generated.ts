import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowSequenceLinkContactMutationVariables = Types.Exact<{
  sequenceId: Types.Scalars['ID']['input'];
  contactId: Types.Scalars['ID']['input'];
  emailId: Types.Scalars['ID']['input'];
}>;

export type FlowSequenceLinkContactMutation = {
  __typename?: 'Mutation';
  flow_sequence_LinkContact: {
    __typename?: 'FlowSequenceContact';
    metadata: { __typename?: 'Metadata'; id: string };
    contact: {
      __typename?: 'Contact';
      metadata: { __typename?: 'Metadata'; id: string };
    };
    email: { __typename?: 'Email'; id: string };
  };
};
