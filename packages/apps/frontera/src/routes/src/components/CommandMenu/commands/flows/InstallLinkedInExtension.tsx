import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const InstallLinkedInExtension = observer(() => {
  const store = useStore();

  const confirmButtonRef = useRef<HTMLButtonElement>(null);

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    window.open(
      'https://chromewebstore.google.com/detail/customeros/khmdccjeodppdldkgifcnkndemjpfoml',
      '_blank',
      'noopener',
    );

    store.ui.commandMenu.toggle('InstallLinkedInExtension');
    store.ui.commandMenu.clearCallback();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold mr-4'>
            Install the LinkedIn browser extension
          </h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <p className='text-sm mt-2'>
          Easily send LinkedIn invitations and messages directly from CustomerOS
          by installing our browser extension
        </p>
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
            Go to extension
          </Button>
        </div>
      </article>
    </Command>
  );
});
