import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const ConfirmBulkFlowEdit = observer(() => {
  const { contacts, ui, flows } = useStore();

  const context = ui.commandMenu.context;

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  const selectedIds = context.ids;

  const selectedFlow = context.property
    ? (flows.value.get(context.property) as FlowStore)
    : null;

  const contactsInFlows = selectedIds
    .map((e) => contacts.value.get(e)?.flowsIds)
    .filter((id) => id !== undefined);

  const handleConfirm = () => {
    if (!context.ids?.length || !context.property) return;

    selectedFlow?.linkContacts(selectedIds);

    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const handleClose = () => {
    ui.commandMenu.toggle('ConfirmBulkFlowEdit');
    ui.commandMenu.clearCallback();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100 cursor-default'>
        <div className='flex justify-between'>
          <h1 className='text-base font-semibold'>
            {contactsInFlows?.length > 1
              ? `${contactsInFlows?.length} selected contacts are already in other flows`
              : `Add ${selectedIds?.length} contacts to flow?`}{' '}
          </h1>
          <div>
            <CommandCancelIconButton onClose={handleClose} />
          </div>
        </div>
        <p className='mt-1 text-sm'>
          Once added, contacts will be scheduled right away if they meet the
          conditions.
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
            data-test='contact-actions-confirm-flow-change'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Add contacts
          </Button>
        </div>
      </article>
    </Command>
  );
});
