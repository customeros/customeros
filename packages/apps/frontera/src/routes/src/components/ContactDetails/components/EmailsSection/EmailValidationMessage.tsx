import { observer } from 'mobx-react-lite';
import { EmailVerificationStatus } from '@finder/components/Columns/contacts/filterTypes';

import { useStore } from '@shared/hooks/useStore';
import { XCircle } from '@ui/media/icons/XCircle.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { HelpCircle } from '@ui/media/icons/HelpCircle.tsx';
import { AlertCircle } from '@ui/media/icons/AlertCircle.tsx';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward.tsx';
import { EmailDeliverable, EmailValidationDetails } from '@graphql/types';
import { CheckCircleBroken } from '@ui/media/icons/CheckCircleBroken.tsx';

interface Props {
  email: string;
  validationDetails: EmailValidationDetails | undefined;
}

export const emailStatuses = {
  DELIVERABLE_NO_RISK: {
    message: 'Deliverable • No risk',
    icon: (
      <CheckCircleBroken className='text-greenLight-500 size-3 first-of-type' />
    ),
    value: EmailVerificationStatus.NoRisk,
  },
  DELIVERABLE_FIREWALL: {
    message: 'Deliverable • Firewall protected',
    icon: <CheckCircleBroken className='text-greenLight-500 size-3' />,
    value: EmailVerificationStatus.FirewallProtected,
  },
  DELIVERABLE_FREE_ACCOUNT: {
    message: 'Deliverable • Free account',
    icon: <CheckCircleBroken className='text-warning-400 size-3' />,
    value: EmailVerificationStatus.FreeAccount,
  },
  CATCH_ALL: {
    message: "Don't know • Catch-all",
    icon: <AlertCircle className='text-gray-500 size-3' />,
    value: EmailVerificationStatus.CatchAll,
  },
  NOT_VERIFIED: {
    message: "Don't know • Not verified yet",
    icon: <HelpCircle className='text-gray-500 size-3' />,
    value: EmailVerificationStatus.NotVerified,
  },
  UNABLE_TO_VALIDATE: {
    message: "Don't know • Unable to validate",
    icon: <HelpCircle className='text-gray-500 size-3' />,
    value: EmailVerificationStatus.UnableToValidate,
  },
  VERIFICATION_IN_PROGRESS: {
    message: "Don't know • Verification in progress",
    icon: <ClockFastForward className='text-primary-600 size-3' />,
    value: EmailVerificationStatus.VerificationInProgress,
  },
  MAILBOX_FULL: {
    message: 'Not deliverable • Mailbox full',
    icon: <XCircle className='text-error-500 size-3' />,
    value: EmailVerificationStatus.MailboxFull,
  },
  INVALID_MAILBOX: {
    message: 'Not deliverable • Mailbox doesn’t exist',
    icon: <XCircle className='text-error-500 size-3' />,
    value: EmailVerificationStatus.InvalidMailbox,
  },
  INCORRECT_FORMAT: {
    message: 'Not deliverable • Incorrect format',
    icon: <XCircle className='text-error-500 size-3' />,
    value: EmailVerificationStatus.IncorrectFormat,
  },
};

function isValidEmail(email: string) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  return emailRegex.test(email);
}

function checkEmailStatus(emailData?: EmailValidationDetails, email?: string) {
  if (!email) {
    return null;
  }

  if (email && !emailData) {
    const isValidSyntax = isValidEmail(email);

    if (!isValidSyntax) return emailStatuses.INCORRECT_FORMAT;
  }

  if (emailData?.deliverable === EmailDeliverable.Deliverable) {
    if (emailData.isFirewalled) return emailStatuses.DELIVERABLE_FIREWALL;
    if (emailData.isFreeAccount) return emailStatuses.DELIVERABLE_FREE_ACCOUNT;

    return emailStatuses.DELIVERABLE_NO_RISK;
  }

  if (emailData?.deliverable === EmailDeliverable.Unknown) {
    return emailData.isCatchAll
      ? emailStatuses.CATCH_ALL
      : emailStatuses.UNABLE_TO_VALIDATE;
  }

  if (emailData?.deliverable === EmailDeliverable.Undeliverable) {
    if (emailData.isMailboxFull) return emailStatuses.MAILBOX_FULL;

    return emailStatuses.INVALID_MAILBOX;
  }

  if (emailData?.verified === false) return emailStatuses.NOT_VERIFIED;
  if (emailData?.verifyingCheckAll)
    return emailStatuses.VERIFICATION_IN_PROGRESS;
}

export const EmailValidationMessage = observer(
  ({ email, validationDetails }: Props) => {
    const data = checkEmailStatus(validationDetails, email);
    const store = useStore();

    if (!data) return null;

    return (
      <Tooltip side='bottom' label={data?.message}>
        <div
          role='button'
          className='flex items-center cursor-pointer'
          onClick={() => {
            store.ui.commandMenu.setType('ContactEmailVerificationInfoModal');
            store.ui.commandMenu.setContext({
              ids: [],
              entity: 'Contact',
              property: data?.value,
            });
            store.ui.commandMenu.setOpen(true);
          }}
        >
          {data?.icon}
        </div>
      </Tooltip>
    );
  },
);
