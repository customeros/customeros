import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

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

    flow?.startFlow({
      onSuccess: () => handleClose(),
    });
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>Start {flow?.value.name}?</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <div className='text-sm mt-2'>
          Making this flow live will trigger it for{' '}
          {flow?.value.participants?.length}{' '}
          {flow?.value.participants?.length === 1 ? 'contact' : 'contacts'}{' '}
          right away and automatically for future contacts when the trigger
          conditions are met.
          <p className='mt-2'>
            We will automatically save your latest changes.
          </p>
        </div>
        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            isLoading={flow?.isLoading}
            data-test='flow-actions-confirm-stop'
            loadingText={`Starting ${flow?.value.name}...`}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Start flow
          </Button>
        </div>
      </article>
    </Command>
  );
});
