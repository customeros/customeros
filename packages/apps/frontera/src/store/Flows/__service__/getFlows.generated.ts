import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetFlowsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type GetFlowsQuery = {
  __typename?: 'Query';
  flows: Array<{
    __typename?: 'Flow';
    name: string;
    description: string;
    status: Types.FlowStatus;
    metadata: { __typename?: 'Metadata'; id: string };
    actions: Array<{
      __typename?: 'FlowAction';
      index: any;
      name: string;
      status: Types.FlowActionStatus;
      actionType: Types.FlowActionType;
      metadata: { __typename?: 'Metadata'; id: string };
      actionData:
        | {
            __typename: 'FlowActionDataEmail';
            replyToId?: string | null;
            subject: string;
            bodyTemplate: string;
          }
        | { __typename: 'FlowActionDataWait'; minutes: any }
        | {
            __typename: 'FlowActionLinkedinConnectionRequest';
            messageTemplate: string;
          }
        | { __typename: 'FlowActionLinkedinMessage'; messageTemplate: string };
      senders: Array<{
        __typename?: 'FlowActionSender';
        mailbox?: string | null;
        metadata: { __typename?: 'Metadata'; id: string };
        user?: {
          __typename?: 'User';
          firstName: string;
          lastName: string;
          name?: string | null;
        } | null;
      }>;
    }>;
    contacts: Array<{
      __typename?: 'FlowContact';
      metadata: { __typename?: 'Metadata'; id: string };
      contact: {
        __typename?: 'Contact';
        metadata: { __typename?: 'Metadata'; id: string };
      };
    }>;
  }>;
};
