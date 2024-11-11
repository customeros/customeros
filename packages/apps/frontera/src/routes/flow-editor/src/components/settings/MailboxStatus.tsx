import { observer } from 'mobx-react-lite';

import { FlowSender } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';

type SenderStatusProps = {
  id: string;
  hasEmailNodes: boolean;
  hasLinkedInNodes: boolean;
};

const EMAILS_PER_MAILBOX = 40;
const LINKEDIN_INVITES_PER_SENDER = 20;

export const SenderStatus = observer(
  ({ id, hasEmailNodes, hasLinkedInNodes }: SenderStatusProps) => {
    const store = useStore();
    const flow = store.flows.value.get(id);
    const senders = flow?.value?.senders ?? [];

    if (!senders.length) {
      return (
        <p className='text-sm'>
          To start sending, add some senders to this flow first
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
    const linkedInCapacity = totalLinkedInSenders * LINKEDIN_INVITES_PER_SENDER;

    // Email node added only, but no mailboxes on any sender
    if (hasEmailNodes && !hasLinkedInNodes && totalMailboxes === 0) {
      return null;
    }

    // LinkedIn node added only, but no extension on any sender
    if (hasLinkedInNodes && !hasEmailNodes && totalLinkedInSenders === 0) {
      return null;
    }

    // Email and LI node added, but no mailboxes or extension yet
    if (
      hasEmailNodes &&
      hasLinkedInNodes &&
      totalMailboxes === 0 &&
      totalLinkedInSenders === 0
    ) {
      return null;
    }

    // Both nodes with both capabilities
    if (
      hasEmailNodes &&
      hasLinkedInNodes &&
      totalMailboxes > 0 &&
      totalLinkedInSenders > 0
    ) {
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

    // Email and LinkedIn nodes with various combinations
    if (hasEmailNodes && hasLinkedInNodes) {
      // Only mailboxes available
      if (totalMailboxes > 0 && totalLinkedInSenders === 0) {
        return (
          <p className='text-sm'>
            You can send up to
            <span className='font-medium mx-1'>{emailCapacity}</span>
            emails across
            <span className='font-medium mx-1'>
              {totalMailboxes} {totalMailboxes === 1 ? 'mailbox' : 'mailboxes'}
            </span>
            mailboxes.
          </p>
        );
      }

      // Only LinkedIn available
      if (totalMailboxes === 0 && totalLinkedInSenders > 0) {
        return (
          <p className='text-sm'>
            You can send up to
            <span className='font-medium mx-1'>{linkedInCapacity}</span>
            <span className='font-medium mx-1'>LinkedIn invites</span>
            per day
          </p>
        );
      }

      // Both available
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
    }

    return null;
  },
);
