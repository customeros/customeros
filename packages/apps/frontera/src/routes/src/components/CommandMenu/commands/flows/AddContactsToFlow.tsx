import React, { useState } from 'react';

import Fuse from 'fuse.js';
import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store';
import { ContactStore } from '@store/Contacts/Contact.store';

import { Avatar } from '@ui/media/Avatar';
import { Check } from '@ui/media/icons/Check';
import { User03 } from '@ui/media/icons/User03';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

export const AddContactsToFlow = observer(() => {
  const { contacts, ui, flows, organizations } = useStore();
  const [search, setSearch] = useState('');

  const context = ui.commandMenu.context;
  const selectedFlowId = context.ids?.[0];

  const selectedFlow = flows.value.get(selectedFlowId) as FlowStore;

  const handleSelect = (opt: ContactStore) => {
    if (!selectedFlowId) {
      ui.toastError('No flow selected', 'no-flow-selected');

      return;
    }
    selectedFlow?.linkContact(opt.id);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const contactsOptions = contacts.toComputedArray((arr) =>
    search
      ? new Fuse(arr, {
          keys: ['value.name'],
          threshold: 0.3,
          isCaseSensitive: false,
        })
          .search(removeAccents(search), { limit: 40 })
          .map((r) => r.item)
      : arr.slice(0, 40),
  );

  return (
    <Command shouldFilter={false} label='Add contact to flow...'>
      <CommandInput
        value={search}
        onValueChange={setSearch}
        placeholder='Add contacts to flow...'
        label={`Flow - ${selectedFlow.value.name}`}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />

      <Command.List>
        {contactsOptions.map((contactStore) => {
          const isSelected = contactStore?.flowsIds?.includes(selectedFlowId);

          return (
            <Command.Item
              key={contactStore.value.metadata.id}
              onSelect={() => {
                handleSelect(contactStore as ContactStore);
              }}
              value={
                contactStore.name || contactStore.value.metadata.id
                  ? `${contactStore.name} ${contactStore.value.metadata.id}`
                  : ''
              }
            >
              <div className='flex justify-between w-full'>
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

                  {contactStore.organizationId && (
                    <span className='ml-1.5 text-gray-500 line-clamp-1 max-w-[250px]'>
                      ·{' '}
                      {organizations.value.get(contactStore.organizationId)
                        ?.value.name ?? ''}
                    </span>
                  )}
                </div>
                {isSelected && <Check />}
              </div>
            </Command.Item>
          );
        })}
      </Command.List>
    </Command>
  );
});

function removeAccents(str: string) {
  return str
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '');
}
