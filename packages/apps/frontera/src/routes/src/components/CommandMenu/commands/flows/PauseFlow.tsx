import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { FlowStatus } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const PauseFlow = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const flow = store.flows.value.get(context.ids?.[0]);

  const confirmButtonRef = useRef<HTMLButtonElement>(null);

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    flow?.update((f) => {
      f.status = FlowStatus.Paused;

      return f;
    });
    handleClose();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>
            Pause flow '{flow?.value.name}'?
          </h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        {/* todo update when we support multiple record types*/}
        <p className='text-sm mt-2'>
          Pausing this flow will immediately stop all upcoming actions for
          active contacts.
          <p className='mt-2'>
            You can resume the flow at any time and these contacts will continue
            from their last completed step.
          </p>
        </p>

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='error'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='flow-actions-confirm-stop'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Pause flow
          </Button>
        </div>
      </article>
    </Command>
  );
});
