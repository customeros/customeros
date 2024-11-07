import React, { useRef, useEffect } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { ContactStore } from '@store/Contacts/Contact.store';
import { FlowContactStore } from '@store/FlowContacts/FlowContact.store.ts';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const UnlinkContactFromFlow = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  const entity = match(context.entity)
    .returnType<ContactStore[] | ContactStore | undefined>()
    .with('Contact', () => store.contacts.value.get(context.ids?.[0]))
    .otherwise(() => undefined);

  const handleClose = () => {
    store.ui.commandMenu.toggle('UnlinkContactFromFlow');
    store.ui.commandMenu.clearCallback();
  };

  const getContactFlowIds = (contactId: string) => {
    const contact = store.contacts.value.get(contactId);

    if (!contact?.flows) return [];

    return contact.flows
      .map((flow) => {
        const matchingContact = flow.value.contacts.find(
          (c) => c.contact.metadata.id === contactId,
        );

        return matchingContact?.metadata.id;
      })
      .filter((id): id is string => typeof id === 'string');
  };
  const flowContactIds: string[] = [
    ...new Set(context.ids.flatMap((id) => getContactFlowIds(id))),
  ];
  const flowContact = store.flowContacts.value.get(
    flowContactIds[0],
  ) as FlowContactStore;

  const handleConfirm = () => {
    if (!context.ids?.length) return;

    if (flowContactIds.length > 1) {
      store.flowContacts.deleteFlowContacts(flowContactIds);
      handleClose();

      return;
    }

    flowContact.deleteFlowContact();
    handleClose();
  };

  const title =
    context.ids?.length > 1
      ? `Remove ${context.ids?.length} contacts from all flows?`
      : `Remove ${(entity as ContactStore)?.name} from ${
          flowContactIds.length === 1
            ? (entity as ContactStore)?.flows?.[0]?.value?.name
            : 'all their flows'
        }?`;

  const description =
    context.ids?.length > 1
      ? `This will remove ${context.ids?.length} contacts from their flows and immediately stop any remaining actions`
      : `This will remove ${flowContact.contact?.name} from their flows and immediately stop any remaining actions`;

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
            colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            dataTest='org-actions-confirm-archive'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Remove {context.ids.length === 1 ? 'contact' : 'contacts'}
          </Button>
        </div>
      </article>
    </Command>
  );
});
