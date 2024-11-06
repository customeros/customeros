import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

export const SenderStatus = observer(
  ({
    id,
    hasEmailNodes,
    hasLinkedInNodes,
  }: {
    id: string;
    hasEmailNodes: boolean;
    hasLinkedInNodes: boolean;
  }) => {
    const store = useStore();
    const flow = store.flows.value.get(id);

    const senders = flow?.value?.senders ?? [];

    const totalMailboxes = senders.reduce((total, sender) => {
      const user = sender?.user?.id && store.users.value.get(sender.user.id);

      if (!user) return total;

      return total + (user?.value?.mailboxes?.length ?? 0);
    }, 0);

    const totalLinkedInSenders = senders.reduce((total, sender) => {
      const user = sender?.user?.id && store.users.value.get(sender.user.id);

      if (!user) return total;

      return total + (user?.value?.hasLinkedInToken ? 1 : 0);
    }, 0);

    const emailCapacity = totalMailboxes * 40;
    const linkedInCapacity = totalLinkedInSenders * 20;

    if (hasEmailNodes && hasLinkedInNodes) {
      if (totalMailboxes > 0 && totalLinkedInSenders > 0) {
        return (
          <p className='text-sm'>
            You can send up to
            <span className='font-medium mx-1'>{emailCapacity}</span>
            emails across
            <span className='font-medium mx-1'>
              {totalMailboxes} {totalMailboxes === 1 ? 'mailbox' : 'mailboxes'}
            </span>
            and a max of
            <span className='font-medium mx-1'>
              {linkedInCapacity} LinkedIn invites
            </span>
            per day.
          </p>
        );
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
            mailboxes. Configure LinkedIn to send invites.
          </p>
        );
      }

      return (
        <p className='text-sm'>
          You haven’t set up any mailboxes yet. Add some to start sending
          emails. You can send up to
          <span className='font-medium mx-1'>
            {linkedInCapacity} LinkedIn invites
          </span>
          per day.
        </p>
      );
    }

    if (hasLinkedInNodes && !hasEmailNodes) {
      return (
        <p className='text-sm'>
          {totalLinkedInSenders > 0
            ? 'Configure LinkedIn to send invites.'
            : `You can send up to 10 LinkedIn invites or 10 messages per day.`}
        </p>
      );
    }

    return (
      <p className='text-sm'>
        You haven’t set up any mailboxes yet. Add some to start sending emails.
      </p>
    );
  },
);
