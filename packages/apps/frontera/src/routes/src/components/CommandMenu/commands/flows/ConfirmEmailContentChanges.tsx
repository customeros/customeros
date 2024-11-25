import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandCancelIconButton } from '@ui/overlay/CommandMenu';

export const ConfirmEmailContentChanges = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const confirmButtonRef = useRef<HTMLButtonElement>(null);

  const handleClose = () => {
    store.ui.commandMenu.clearContext();
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleCancelChanges = () => {
    handleClose();
    store.ui.flowActionSidePanel.setOpen(false);
    store.ui.flowActionSidePanel.clearContext();
  };

  const handleConfirm = () => {
    context?.callback?.();
    store.ui.commandMenu.clearContext();
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>
            Want to save this email first?
          </h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <p className='text-sm mt-2'>
          This email has unsaved changes. Would you like to save them first
          before heading out?
        </p>
        <div className='flex justify-between gap-3 mt-6'>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            onClick={handleCancelChanges}
            data-test={'email-content-leave-without-saving'}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.stopPropagation();
                handleCancelChanges();
              }
            }}
          >
            Discard changes
          </Button>

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test={'email-content-leave-and-save'}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Save changes
          </Button>
        </div>
      </article>
    </Command>
  );
});
