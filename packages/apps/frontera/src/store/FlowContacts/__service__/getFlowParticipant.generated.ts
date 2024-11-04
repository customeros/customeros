import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetFlowParticipantQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type GetFlowParticipantQuery = {
  __typename?: 'Query';
  flowParticipant: {
    __typename?: 'FlowContact';
    status: Types.FlowParticipantStatus;
    scheduledAction?: string | null;
    scheduledAt?: any | null;
    metadata: { __typename?: 'Metadata'; id: string };
    contact: {
      __typename?: 'Contact';
      metadata: { __typename?: 'Metadata'; id: string };
    };
  };
};
