import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const ConfirmSingleFlowEdit = observer(() => {
  const { contacts, ui, flows } = useStore();

  const context = ui.commandMenu.context;

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  const selectedIds = context.ids;

  const selectedFlow = context.property
    ? (flows.value.get(context.property) as FlowStore)
    : null;

  const contact = contacts.value.get(selectedIds[0]);
  const contactsInFlows = contact?.flows?.length ?? 0;

  const handleConfirm = () => {
    if (!context.ids?.length || !context.property) return;

    selectedFlow?.linkContact(selectedIds[0]);

    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const handleClose = () => {
    ui.commandMenu.toggle('ConfirmSingleFlowEdit');
    ui.commandMenu.clearCallback();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100 cursor-default'>
        <div className='flex justify-between'>
          <h1 className='text-base font-semibold'>
            {contactsInFlows > 0
              ? `${contact?.name} is already in other flows`
              : `Add ${contact?.name} to flow?`}
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
            Add contact
          </Button>
        </div>
      </article>
    </Command>
  );
});
