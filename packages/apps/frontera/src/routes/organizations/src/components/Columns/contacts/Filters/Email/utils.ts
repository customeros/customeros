export enum DeliverabilityStatus {
  Deliverable = 'deliverable',
  NotDeliverable = 'not_deliverable',
  Unknown = 'unknown',
}

export const CategoryHeaderLabel = {
  [DeliverabilityStatus.Deliverable]: 'Deliverable',
  [DeliverabilityStatus.NotDeliverable]: 'Not Deliverable',
  [DeliverabilityStatus.Unknown]: 'Don’t know',
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

  // category
  IsDeliverable = 'is_deliverable',
  IsNotDeliverable = 'is_not_deliverable',
  IsUnknown = 'is_unknown',
}

export const getCategoryString = (
  category: DeliverabilityStatus,
): EmailVerificationStatus => {
  switch (category) {
    case DeliverabilityStatus.Deliverable:
      return EmailVerificationStatus.IsDeliverable;
    case DeliverabilityStatus.NotDeliverable:
      return EmailVerificationStatus.IsNotDeliverable;
    case DeliverabilityStatus.Unknown:
      return EmailVerificationStatus.IsUnknown;
  }
};

export const getOptionsForCategory = (category: DeliverabilityStatus) => {
  switch (category) {
    case DeliverabilityStatus.Deliverable:
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
    case DeliverabilityStatus.NotDeliverable:
      return [
        {
          label: 'Invalid mailbox',
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
    case DeliverabilityStatus.Unknown:
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
