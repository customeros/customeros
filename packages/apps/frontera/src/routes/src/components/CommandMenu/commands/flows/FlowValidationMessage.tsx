import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandCancelIconButton } from '@ui/overlay/CommandMenu';

export const FlowValidationMessage = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const confirmButtonRef = useRef<HTMLButtonElement>(null);

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setContext({
      ...context,
      meta: undefined,
    });
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    context.callback?.();
    handleClose();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>{context.meta?.title}</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <p className='text-sm mt-2'>{context.meta?.description}</p>
        <div className='flex justify-between gap-3 mt-6'>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='flow-validation-message-close-button'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            {context.meta?.buttonText}
          </Button>
        </div>
      </article>
    </Command>
  );
});
