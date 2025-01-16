import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const StopFlow = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const flow = store.flows.value.get(context.ids?.[0]);

  const confirmButtonRef = useRef<HTMLButtonElement>(null);

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    flow?.stopFlow({
      onSuccess: () => handleClose(),
    });
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>Stop {flow?.value.name}?</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        {/* todo update when we support multiple record types*/}
        <div className='text-sm mt-2'>
          This will stop all upcoming steps for your active contacts from taking
          place.
          <p className='mt-2'>
            When you start the flow again, contacts will pick up from the last
            completed step on a new schedule.
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
            loadingText={`Stopping ${flow?.value.name}...`}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Stop flow
          </Button>
        </div>
      </article>
    </Command>
  );
});
