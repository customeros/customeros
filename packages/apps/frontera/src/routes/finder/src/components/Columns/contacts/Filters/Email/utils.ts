import { EmailDeliverable } from '@graphql/types';

export const CategoryHeaderLabel = {
  [EmailDeliverable.Deliverable]: 'Deliverable',
  [EmailDeliverable.Undeliverable]: 'Not Deliverable',
  [EmailDeliverable.Unknown]: 'Don’t know',
};

export enum EmailVerificationStatus {
  NoRisk = 'no_risk',
  FirewallProtected = 'firewall_protected',
  FreeAccount = 'free_account',
  GroupMailbox = 'group_mailbox',
  InvalidMailbox = 'invalid_mailbox',
  MailboxFull = 'mailbox_full',
  IncorrectFormat = 'incorrect_format',
  CatchAll = 'catch_all',
  NotVerified = 'not_verified',
  VerificationInProgress = 'verification_in_progress',
}

export const getOptionsForCategory = (category: EmailDeliverable) => {
  switch (category) {
    case EmailDeliverable.Deliverable:
      return [
        { label: 'No risk', value: EmailVerificationStatus.NoRisk },
        {
          label: 'Firewall protected',
          value: EmailVerificationStatus.FirewallProtected,
        },
        {
          label: 'Free account',
          value: EmailVerificationStatus.FreeAccount,
        },
        {
          disabled: true,
          label: 'Group mailbox',
          value: EmailVerificationStatus.GroupMailbox,
        },
      ];
    case EmailDeliverable.Undeliverable:
      return [
        {
          label: 'Mailbox doesn’t exist',
          value: EmailVerificationStatus.InvalidMailbox,
        },
        {
          label: 'Mailbox full',
          value: EmailVerificationStatus.MailboxFull,
        },
        {
          label: 'Incorrect email format',
          value: EmailVerificationStatus.IncorrectFormat,
        },
      ];
    case EmailDeliverable.Unknown:
      return [
        { label: 'Catch all', value: EmailVerificationStatus.CatchAll },
        {
          label: 'Not verified yet',
          value: EmailVerificationStatus.NotVerified,
        },
        {
          label: 'Verification in progress',
          value: EmailVerificationStatus.VerificationInProgress,
        },
      ];
    default:
      return [];
  }
};
