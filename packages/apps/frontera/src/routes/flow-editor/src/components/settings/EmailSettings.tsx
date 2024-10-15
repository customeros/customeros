import React from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Divider } from '@ui/presentation/Divider';
import { FlowSender as TFlowSender } from '@graphql/types';

import { FlowSender } from './FlowSender';
import { SenderDropdown } from './SenderDropdown';
import { MailboxStatus } from './MailboxStatus.tsx';

export const EmailSettings = observer(({ id }: { id: string }) => {
  const store = useStore();
  const flow = store.flows.value.get(id);

  const hasSenders =
    !!flow?.value?.senders?.length && flow?.value?.senders?.length > 0;

  const totalMailboxes = getTotalMailboxes(flow?.value.senders ?? []);

  return (
    <>
      <h2 className='text-base font-medium'>Flow Settings</h2>
      <div className='flex flex-col gap-2'>
        <h3 className='text-sm font-medium'>Senders</h3>
        <MailboxStatus
          hasSenders={hasSenders}
          totalMailboxes={totalMailboxes}
        />
        <div className='flex flex-col gap-2 '>
          {hasSenders &&
            flow?.value.senders.map((e) => (
              <FlowSender flowId={id} id={e.metadata.id} key={e.metadata.id} />
            ))}
        </div>
        <SenderDropdown flowId={id} />
      </div>
      <Divider />
      <div className='flex flex-col gap-1'>
        <h3 className='text-sm font-medium'>Send schedule</h3>
        <p className='text-sm'>
          Messages will be sent on week days between 8am and 7pm using the
          contact’s timezone
        </p>
      </div>
      <Divider />
      <div className='flex flex-col gap-1'>
        <h3 className='text-sm font-medium'>Opt-out link</h3>
        <p className='text-sm'>
          Opt-out links let recipients easily unsubscribe, helping prevent spam
          reports and keeping communication respectful
        </p>
        <div className='text-sm border px-2 py-2 rounded-md mt-2'>
          If I missed the mark, please{' '}
          <span className='underline'>let me know</span>
        </div>
      </div>
    </>
  );
});

function getTotalMailboxes(senders: TFlowSender[]): number {
  return senders.reduce((total, sender) => {
    if (sender.user) {
      return total + sender.user.mailboxes.length;
    }

    return total;
  }, 0);
}
