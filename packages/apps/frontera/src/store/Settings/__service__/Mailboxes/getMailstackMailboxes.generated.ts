import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type MailstackMailboxesQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type MailstackMailboxesQuery = {
  __typename?: 'Query';
  mailstack_Mailboxes: Array<{
    __typename?: 'Mailbox';
    domain: string;
    mailbox: string;
    userId?: string | null;
    rampUpRate: number;
    rampUpMax: number;
    rampUpCurrent: number;
  }>;
};
