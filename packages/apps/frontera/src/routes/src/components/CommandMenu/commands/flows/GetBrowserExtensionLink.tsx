import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Command, CommandCancelIconButton } from '@ui/overlay/CommandMenu';

export const GetBrowserExtensionLink = observer(() => {
  const store = useStore();

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const [_, copyToClipboard] = useCopyToClipboard();

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    copyToClipboard(
      'https://chromewebstore.google.com/detail/customeros/khmdccjeodppdldkgifcnkndemjpfoml',
      'Link copied to clipboard',
    );

    store.ui.commandMenu.toggle('GetBrowserExtensionLink');
    store.ui.commandMenu.clearCallback();
  };

  const user = store.users.value.get(store.ui.commandMenu.context.ids[0]);

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>
            Ask {user?.name} to install our LinkedIn browser extension
          </h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <p className='text-sm mt-2'>
          For {user?.name} to send LinkedIn invitations and messages, they’d
          need to sign into CustomerOS first and install our browser extension.
        </p>
        <div className='flex justify-between gap-3 mt-6'>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            // colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='flow-actions-confirm-stop'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Copy link to extension
          </Button>
        </div>
      </article>
    </Command>
  );
});
