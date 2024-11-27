import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetRegisteredBuyDomainsWithMailboxesQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type GetRegisteredBuyDomainsWithMailboxesQuery = {
  __typename?: 'Query';
  mailstack_RegisteredBuyDomainsWithMailboxes: Array<{
    __typename?: 'MailstackBuyRequest';
    id: string;
    createdAt: any;
    domains: Array<{
      __typename?: 'MailstackBuyRequestDomain';
      domain: string;
      status: Types.MailstackBuyRequestDomainStatus;
    }>;
    mailboxes: Array<{
      __typename?: 'MailstackBuyRequestMailbox';
      mailbox: string;
      status: Types.MailstackBuyRequestMailboxStatus;
    }>;
  }>;
};
