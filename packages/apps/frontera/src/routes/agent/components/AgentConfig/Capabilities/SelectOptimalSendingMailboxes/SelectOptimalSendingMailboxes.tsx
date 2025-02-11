import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { FlowSender as FlowSenderT } from '@graphql/types';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';

import { FlowSender, MailboxStatus, SenderDropdown } from './components';

export const SelectOptimalSendingMailboxes = observer(() => {
  const store = useStore();
  const flow = store.flows.toArray()?.[0]; // todo get proper flow when BE is ready

  const id = flow?.value?.metadata?.id;
  const hasSenders =
    !!flow?.value?.senders?.length && flow?.value?.senders?.length > 0;

  return (
    <article className='flex flex-col gap-4 -mr-4'>
      <div className='flex flex-col gap-1'>
        <h1 className='text-sm font-medium pr-4'>
          Select optimal sending mailboxes
        </h1>
        <p className='text-sm'>
          We will select the optimal mailbox for sending based on its warming
          status and schedule
        </p>
      </div>

      <ScrollAreaRoot>
        <ScrollAreaViewport className=''>
          <div className=' h-[90vh] pr-4'>
            <div className='flex flex-col gap-2'>
              <h3 className='text-sm font-medium'>Senders</h3>
              <MailboxStatus id={id ?? ''} />
              {hasSenders && (
                <div className='flex flex-col gap-2 '>
                  {flow?.value.senders.map((e: FlowSenderT) => (
                    <FlowSender
                      flowId={id ?? ''}
                      id={e.metadata.id}
                      key={e.metadata.id}
                    />
                  ))}
                </div>
              )}
              <SenderDropdown flowId={id ?? ''} />
            </div>
          </div>
        </ScrollAreaViewport>
        <ScrollAreaScrollbar orientation='vertical'>
          <ScrollAreaThumb />
        </ScrollAreaScrollbar>
      </ScrollAreaRoot>
    </article>
  );
});
