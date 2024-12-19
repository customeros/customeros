import type { Contact } from '@store/Contacts/Contact.dto';

import { useRef, useEffect } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { FlowParticipantStore } from '@store/FlowParticipants/FlowParticipant.store';

import { TableIdType } from '@graphql/types';
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
    .returnType<Contact | undefined>()
    .with(
      'Contact',
      () => store.contacts.value.get(context.ids?.[0] as string) as Contact,
    )
    .otherwise(() => undefined);

  const handleClose = () => {
    store.ui.commandMenu.toggle('UnlinkContactFromFlow');
    store.ui.commandMenu.clearCallback();
  };

  const getContactFlowIds = (contactId: string) => {
    const contact = store.contacts.value.get(contactId);

    return match(context.meta?.tableId)
      .with(TableIdType.FlowContacts, () => {
        const flow = store.flows.value.get(context.meta?.id)?.value
          .participants;

        if (!flow) return [];

        const flowContactId = flow.find((p) => p.entityId === contactId)
          ?.metadata.id;

        return flowContactId ? [flowContactId] : [];
      })
      .otherwise(() => {
        if (!contact?.flows) return [];

        return contact.flows
          ?.map((flow) => {
            const matchingContact = flow.value.participants.find(
              (c) => c.entityId === contactId,
            );

            return matchingContact?.metadata.id;
          })
          .filter((id): id is string => typeof id === 'string');
      });
  };
  const flowContactIds: string[] = [
    ...new Set(context.ids.flatMap((id) => getContactFlowIds(id))),
  ];

  const flowContact = store.flowParticipants.value.get(
    flowContactIds[0],
  ) as FlowParticipantStore;

  const handleConfirm = () => {
    if (!context.ids?.length) return;

    if (flowContactIds.length > 1) {
      store.flowParticipants.deleteFlowParticipants(
        flowContactIds,
        context?.meta?.id,
      );
      handleClose();

      return;
    }

    store.flowParticipants.deleteFlowParticipant(
      context?.meta?.id,
      flowContactIds[0],
    );
    handleClose();
  };

  const title = match(context.meta?.tableId)
    .with(TableIdType.FlowContacts, () => {
      const flowName = store.flows.value.get(context.meta?.id)?.value.name;

      return context.ids?.length > 1
        ? `Remove ${context.ids?.length} contacts from ${flowName}?`
        : `Remove ${(entity as Contact)?.name} from ${flowName}?`;
    })
    .otherwise(() => {
      return context.ids?.length > 1
        ? `Remove ${context.ids?.length} contacts from all flows?`
        : `Remove ${(entity as Contact)?.name} from ${
            flowContactIds.length === 1
              ? (entity as Contact)?.flows?.[0]?.value?.name
              : 'all their flows'
          }?`;
    });

  const description = match(context.meta?.tableId)
    .with(TableIdType.FlowContacts, () => {
      return context.ids?.length > 1
        ? `This will remove ${context.ids?.length} contacts from this flow and immediately stop any remaining actions`
        : `This will remove ${flowContact.contact?.name} from this flow and immediately stop any remaining actions`;
    })
    .otherwise(() => {
      return context.ids?.length > 1
        ? `This will remove ${context.ids?.length} contacts from their flows and immediately stop any remaining actions`
        : `This will remove ${flowContact.contact?.name} from their flows and immediately stop any remaining actions`;
    });

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
