import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

export const MailboxStatus = observer(({ id }: { id: string }) => {
  const store = useStore();
  const flow = store.flows.value.get(id);

  const hasSenders =
    !!flow?.value?.senders?.length && flow?.value?.senders?.length > 0;

  if (!hasSenders) {
    return (
      <p className='text-sm'>
        Add one or more senders to start sending emails in this flow
      </p>
    );
  }

  const totalMailboxes = flow?.value.senders.reduce((total, sender) => {
    const user = sender?.user?.id && store.users.value.get(sender.user.id);

    if (user) {
      return total + (user?.value?.mailboxes?.length ?? 0);
    }

    return total;
  }, 0);

  return (
    <p className='text-sm'>
      You have{' '}
      <span className='font-medium'>
        {totalMailboxes} {totalMailboxes === 1 ? 'mailbox' : 'mailboxes'}
      </span>{' '}
      available allowing you to send up to{' '}
      <span className='font-medium'>{totalMailboxes * 40} emails per day.</span>
    </p>
  );
});
