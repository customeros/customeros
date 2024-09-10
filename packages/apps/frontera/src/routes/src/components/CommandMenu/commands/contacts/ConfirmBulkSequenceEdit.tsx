import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const ConfirmBulkSequenceEdit = observer(() => {
  const { contacts, ui, flowSequences } = useStore();

  const context = ui.commandMenu.context;

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  const selectedIds = context.ids;

  const selectedSequence = context.property
    ? flowSequences.value.get(context.property)
    : null;

  const contactsInSequences = selectedIds
    .map((e) => contacts.value.get(e)?.sequence?.value?.name)
    .filter((name) => name !== undefined);

  const handleConfirm = () => {
    if (!context.ids?.length || !context.property) return;

    flowSequences.linkContacts(context.property, selectedIds);

    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const handleClose = () => {
    ui.commandMenu.toggle('ConfirmBulkSequenceEdit');
    ui.commandMenu.clearCallback();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>
            {contactsInSequences?.length} of your selected contacts are already
            in other sequences
          </h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>
        <p className='mt-1 text-sm'>
          To add them to {selectedSequence?.value?.name}, we’ll have to remove
          them from their existing sequences.
        </p>

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton ref={closeButtonRef} onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='contact-actions-confirm-sequence-change'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Continue
          </Button>
        </div>
      </article>
    </Command>
  );
});
