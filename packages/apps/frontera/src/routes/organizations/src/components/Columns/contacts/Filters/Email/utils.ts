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
}
