import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetRegisteredBuyDomainsWithMailboxesQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type GetRegisteredBuyDomainsWithMailboxesQuery = {
  __typename?: 'Query';
  mailstack_RegisteredBuyDomainsWithMailboxes: Array<{
    __typename?: 'RegisteredBuyDomainWithMailboxes';
    id: string;
    createdAt: any;
    domain: {
      __typename?: 'MailstackDomain';
      domain: string;
      status: Types.MailstackStatus;
    };
    mailboxes: Array<{
      __typename?: 'MailstackMailbox';
      mailbox: string;
      status: Types.MailstackStatus;
    }>;
  }>;
};
