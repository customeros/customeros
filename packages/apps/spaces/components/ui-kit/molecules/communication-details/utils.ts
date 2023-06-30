export interface MappedObject {
  [key: string]: {
    condition: boolean | Array<string>;
    message: string;
  };
}

export interface InputObject {
  [key: string]: boolean | null | string;
}

export const VALIDATION_MESSAGES: MappedObject = {
  isReachable: {
    condition: ['invalid', 'risky', 'FALSE'],
    message: 'This mailbox is not reachable.',
  },
  isValidSyntax: {
    condition: false,
    message: 'Not a valid email address. Check for typos?',
  },
  canConnectSmtp: {
    condition: false,
    message: '',
  },
  acceptsMail: {
    condition: false,
    message: 'This domain does not accept emails.',
  },
  hasFullInbox: {
    condition: true,
    message: 'This mailbox is full. Your message is likely to bounce.',
  },
  isCatchAll: {
    condition: true,
    message:
      'This catch-all mailbox might not belong to your intended recipient.',
  },
  isDeliverable: {
    condition: false,
    message: 'This mailbox does not accept emails.',
  },
  isDisabled: {
    condition: true,
    message: 'This mailbox is disabled.',
  },
};
