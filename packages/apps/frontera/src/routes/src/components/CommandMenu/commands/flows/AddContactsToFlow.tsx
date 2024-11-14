import { useState } from 'react';

import Fuse from 'fuse.js';
import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store';
import { ContactStore } from '@store/Contacts/Contact.store.ts';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const AddContactsToFlow = observer(() => {
  const { contacts, ui, flows } = useStore();
  const [search, setSearch] = useState('');

  const context = ui.commandMenu.context;
  const selectedFlowId = context.ids?.[0];

  const selectedFlow = flows.value.get(selectedFlowId) as FlowStore;

  const handleSelect = (opt: ContactStore) => {
    if (!selectedFlowId) {
      ui.toastError('No flow selected', 'no-flow-selected');
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
        placeholder='Search contacts...'
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
            <CommandItem
              key={contactStore.id}
              rightAccessory={isSelected ? <Check /> : undefined}
              onSelect={() => {
                handleSelect(contactStore as ContactStore);
              }}
            >
              {contactStore.name?.trim()?.length
                ? contactStore.name
                : 'Unnamed'}
            </CommandItem>
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
