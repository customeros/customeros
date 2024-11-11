import React, { useRef } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { FlowStatus } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const StartFlow = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const flow = store.flows.value.get(context.ids?.[0]);

  const confirmButtonRef = useRef<HTMLButtonElement>(null);

  const handleClose = () => {
    store.ui.commandMenu.clearContext();
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    if (context?.meta?.hasUnsavedChanges) {
      context?.callback?.();
      handleClose();

      return;
    }

    flow?.update((f) => {
      f.status = FlowStatus.Scheduling;

      return f;
    });
    store.ui.commandMenu.setOpen(false);

    store.ui.commandMenu.clearContext();
  };

  const data = match(flow?.value?.status)
    .with(FlowStatus.Paused, () => ({
      title: `Resume flow '${flow?.value.name}'?`,
      description: `Resuming this flow will immediately start all upcoming actions for active contacts.`,
      button: 'Resume flow',
    }))
    .otherwise(() => ({
      title: `Start flow '${flow?.value.name}'?`,
      description: (
        <>
          Making this flow live will trigger it for{' '}
          {flow?.value.contacts?.length}{' '}
          {flow?.value.contacts?.length === 1 ? 'contact' : 'contacts'} right
          away and automatically for future contacts when the trigger conditions
          are met.
          <p className='mt-2'>
            We will automatically save your latest changes.
          </p>
        </>
      ),
      button: 'Start flow',
    }));

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>{data.title}</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <p className='text-sm mt-2'>{data.description}</p>
        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='flow-actions-confirm-stop'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            {data.button}
          </Button>
        </div>
      </article>
    </Command>
  );
});
