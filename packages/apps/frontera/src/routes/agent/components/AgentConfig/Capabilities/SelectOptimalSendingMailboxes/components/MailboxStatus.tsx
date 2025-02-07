import { observer } from 'mobx-react-lite';

import { FlowSender } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';

type MailboxStatusProps = {
  id: string;
};

const EMAILS_PER_MAILBOX = 40;

export const MailboxStatus = observer(({ id }: MailboxStatusProps) => {
  const store = useStore();
  const flow = store.flows.value.get(id);
  const senders = flow?.value?.senders ?? [];

  if (!senders.length) {
    return (
      <p className='text-sm'>
        To start sending, add some senders to this agent first
      </p>
    );
  }

  const totalMailboxes = senders.reduce((total, sender: FlowSender) => {
    const user = sender?.user?.id && store.users.value.get(sender.user.id);

    if (!user) return total;

    return total + (user?.value?.mailboxes?.length ?? 0);
  }, 0);

  const totalLinkedInSenders = senders.reduce((total, sender: FlowSender) => {
    const user = sender?.user?.id && store.users.value.get(sender.user.id);

    if (!user) return total;

    return total + (user?.value?.hasLinkedInToken ? 1 : 0);
  }, 0);

  const emailCapacity = totalMailboxes * EMAILS_PER_MAILBOX;

  if (totalMailboxes === 0 && totalLinkedInSenders === 0) {
    return null;
  }

  if (totalMailboxes > 0) {
    return (
      <p className='text-sm'>
        You can send up to
        <span className='font-medium mx-1'>{emailCapacity}</span>
        emails across
        <span className='font-medium mx-1'>
          {totalMailboxes} {totalMailboxes === 1 ? 'mailbox' : 'mailboxes'}
        </span>
      </p>
    );
  }

  return null;
});
