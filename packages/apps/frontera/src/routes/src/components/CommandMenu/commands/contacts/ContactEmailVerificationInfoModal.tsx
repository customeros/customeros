import React, { useRef, MouseEvent, KeyboardEvent } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { EmailVerificationStatus } from '@finder/components/Columns/contacts/filterTypes';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandCancelIconButton } from '@ui/overlay/CommandMenu';

export const ContactEmailVerificationInfoModal = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const closeButtonRef = useRef<HTMLButtonElement>(null);

  const handleClose = (
    e: MouseEvent<HTMLButtonElement> | KeyboardEvent<HTMLButtonElement>,
  ) => {
    e.stopPropagation();
    e.preventDefault();
    store.ui.commandMenu.toggle('ContactEmailVerificationInfoModal');
    store.ui.commandMenu.clearContext();
  };

  const data = match(context.property)
    .with(EmailVerificationStatus.NoRisk, () => ({
      title: 'Deliverable • No risk',
      description: (
        <>
          <p>
            These email addresses are verified and safe to send messages to.
          </p>
          <p>
            You can send to these with confidence, knowing that they won't be
            blocked or result in bounces.
          </p>
        </>
      ),
    }))
    .with(EmailVerificationStatus.FirewallProtected, () => ({
      title: 'Deliverable • Firewall protected',
      description: (
        <>
          <p>
            These email addresses are deliverable, but the recipient's server is
            protected by a firewall, which may block or filter some emails.
          </p>
          <p>
            We recommend sending to these from your personal mailbox or from
            outbound mailboxes that have a strong reputation and have been
            active for over 30 days.
          </p>
        </>
      ),
    }))
    .with(EmailVerificationStatus.FreeAccount, () => ({
      title: 'Deliverable • Free account',
      description: (
        <>
          <p>
            These email addresses are associated with a free email service
            provider, such as Gmail or Yahoo and have a higher risk of being
            "burner" accounts.
          </p>
          <p>
            While the email is deliverable, engagement will be lower as we don't
            know if these are actively monitored.
          </p>
        </>
      ),
    }))
    .with(EmailVerificationStatus.GroupMailbox, () => ({
      title: 'Deliverable • Group mailbox',
      description: (
        <>
          <p>
            These email addresses are associated with a group or distribution
            list, meaning emails will be delivered to multiple recipients.
            Examples include sales@ or support@.
          </p>
          <p>We automatically hide these mailboxes from views.</p>
        </>
      ),
    }))
    .with(EmailVerificationStatus.MailboxFull, () => ({
      title: 'Not deliverable • Mailbox full',
      description: (
        <>
          <p>
            These email addresses are currently undeliverable because the
            recipient's mailbox is full and cannot receive new emails.
          </p>
          <p>You can retry sending to them at a later stage.</p>
        </>
      ),
    }))
    .with(EmailVerificationStatus.IncorrectFormat, () => ({
      title: 'Not deliverable • Incorrect format',
      description: (
        <>
          <p>
            These email addresses are undeliverable because they are formatted
            incorrectly.
          </p>
          <p>
            Please check for any missing characters or symbols. Once the format
            is corrected, we'll automatically retry verifying them.
          </p>
        </>
      ),
    }))
    .with(EmailVerificationStatus.InvalidMailbox, () => ({
      title: `Not deliverable • Mailbox doesn't exist`,
      description: (
        <p>
          These email addresses are undeliverable because the mailbox doesn't
          exist.
        </p>
      ),
    }))
    .with(EmailVerificationStatus.CatchAll, () => ({
      title: `Don't know • Catch-all`,
      description: (
        <>
          <p>
            Catch-all domains accept all emails sent to them, even if the
            specific address doesn't exist, which can make deliverability
            uncertain.
          </p>
          <p>
            More advanced checks are needed to verify if emails will actually be
            delivered.
          </p>
        </>
      ),
    }))
    .with(EmailVerificationStatus.NotVerified, () => ({
      title: `Don't know • Not verified yet`,
      description: (
        <p>These email addresses have not been verified by CustomerOS yet.</p>
      ),
    }))
    .with(EmailVerificationStatus.VerificationInProgress, () => ({
      title: `Don't know • Verification in progress`,
      description: (
        <p>
          These emails are currently being verified. This typically takes one or
          two days.
        </p>
      ),
    }))
    .otherwise(() => ({
      title: '',
      description: '',
    }));

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between mb-0.5'>
          <h1 className='text-base font-semibold'>{data.title}</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <div className='text-sm flex flex-col gap-4'>{data.description}</div>

        <div className='flex justify-between gap-3 mt-6'>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            ref={closeButtonRef}
            onClick={handleClose}
            onFocus={(e) => e.preventDefault()}
          >
            Got it
          </Button>
        </div>
      </article>
    </Command>
  );
});
