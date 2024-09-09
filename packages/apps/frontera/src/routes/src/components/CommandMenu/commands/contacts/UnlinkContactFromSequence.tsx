import React, { useRef, useEffect } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { ContactStore } from '@store/Contacts/Contact.store';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const UnlinkContactFromSequence = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  const entity = match(context.entity)
    .returnType<ContactStore[] | ContactStore | undefined>()
    .with('Contact', () => store.contacts.value.get(context.ids?.[0]))
    .otherwise(() => undefined);

  const handleClose = () => {
    store.ui.commandMenu.toggle('DeleteConfirmationModal');
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    handleClose();
  };

  const title =
    context.ids?.length > 1
      ? `Archive ${context.ids?.length} contacts?`
      : `Archive ${(entity as ContactStore)?.value.name}?`;

  const description =
    context.ids?.length > 1
      ? `Removing ${context.ids?.length} contacts will end their sequences`
      : `Removing ${
          (entity as ContactStore)?.value?.name
        } will end the sequence for them`;

  useEffect(() => {
    closeButtonRef.current?.focus();
  }, []);

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>{title}</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>
        {description && <p className='mt-1 text-sm'>{description}</p>}

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton ref={closeButtonRef} onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='error'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='org-actions-confirm-archive'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Remove contacts
          </Button>
        </div>
      </article>
    </Command>
  );
});
