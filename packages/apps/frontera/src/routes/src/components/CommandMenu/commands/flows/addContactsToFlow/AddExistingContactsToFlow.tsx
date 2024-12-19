import React, { useState } from 'react';

import Fuse from 'fuse.js';
import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store';
import { Contact } from '@store/Contacts/Contact.dto';
import { FlowParticipantStore } from '@store/FlowParticipants/FlowParticipant.store.ts';

import { Avatar } from '@ui/media/Avatar';
import { Check } from '@ui/media/icons/Check';
import { User03 } from '@ui/media/icons/User03';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

export const AddExistingContacts = observer(() => {
  const { contacts, ui, flows, flowParticipants } = useStore();
  const [search, setSearch] = useState('');

  const context = ui.commandMenu.context;
  const selectedFlowId = context.ids?.[0];

  const selectedFlow = flows.value.get(selectedFlowId) as FlowStore;

  const handleSelect = (opt: Contact) => {
    if (!selectedFlowId) {
      ui.toastError('No flow selected', 'no-flow-selected');

      return;
    }
    const isSelected = opt?.flowsIds?.includes(selectedFlowId);

    if (isSelected) {
      const participant = opt?.flows
        ?.find((flow) => flow?.id === selectedFlowId)
        ?.value.participants.find(
          (participant) => participant.entityId === opt.id,
        );

      const flowParticipant =
        participant?.metadata.id &&
        (flowParticipants.value.get(
          participant.metadata.id,
        ) as FlowParticipantStore);

      if (flowParticipant) {
        flowParticipants?.deleteFlowParticipant(
          selectedFlowId,
          flowParticipant.id,
        );

        return;
      }
    } else {
      selectedFlow?.linkContact(opt.id);
    }
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const contactsOptions = contacts.toComputedArray((arr) =>
    search
      ? new Fuse(arr, {
          keys: ['value.name', 'value.firstName', 'value.lastName'],
          threshold: 0.3,
          isCaseSensitive: false,
        })
          .search(removeAccents(search), { limit: 10 })
          .map((r) => r.item)
      : arr.slice(0, 10),
  );

  return (
    <>
      <div className='-mt-3 '>
        <CommandInput
          value={search}
          onValueChange={setSearch}
          className='text-sm p-0 -ml-2'
          placeholder='Search a contact…'
          onKeyDownCapture={(e) => {
            if (e.key === ' ') {
              e.stopPropagation();
            }
          }}
        />
      </div>

      <Command.List className='!px-0'>
        {contactsOptions.map((contactStore) => {
          const isSelected = contactStore?.flowsIds?.includes(selectedFlowId);

          return (
            <Command.Item
              key={contactStore.value.id}
              onSelect={() => {
                handleSelect(contactStore as Contact);
              }}
              value={
                contactStore.name || contactStore.value.id
                  ? `${contactStore.name} ${contactStore.value.id}`
                  : ''
              }
            >
              <div className='flex justify-between w-full items-center'>
                <div className='flex items-center'>
                  <Avatar
                    size='xxs'
                    textSize='xxs'
                    name={contactStore.name}
                    icon={<User03 className='text-primary-700  ' />}
                    src={
                      contactStore?.value?.profilePhotoUrl
                        ? contactStore.value.profilePhotoUrl
                        : undefined
                    }
                  />
                  <span className='ml-2 capitalize line-clamp-1 '>
                    {contactStore.name?.length ? contactStore.name : 'Unnamed'}
                  </span>

                  <span className='ml-1.5 text-gray-500 line-clamp-1 max-w-[250px]'>
                    ·{' '}
                    {contactStore.value.emails.find((e) => e.primary)?.email ??
                      'No email yet'}
                  </span>
                </div>
                {isSelected && <Check />}
              </div>
            </Command.Item>
          );
        })}
      </Command.List>
    </>
  );
});

function removeAccents(str: string) {
  return str
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '');
}
