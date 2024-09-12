import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { ContactStore } from '@store/Contacts/Contact.store';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store';

import { Plus } from '@ui/media/icons/Plus';
import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditContactSequence = observer(() => {
  const { contacts, ui, flowSequences } = useStore();
  const [search, setSearch] = useState('');

  const context = ui.commandMenu.context;

  const contact = contacts.value.get(context.ids?.[0] as string);
  const selectedIds = context.ids;

  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.value.name}`
      : `${selectedIds?.length} contacts`;

  const handleOpenConfirmDialog = (id: string) => {
    ui.commandMenu.toggle('ConfirmBulkSequenceEdit');
    ui.commandMenu.setContext({
      ...ui.commandMenu.context,
      property: id,
    });
    ui.commandMenu.setOpen(true);
  };

  const handleSelect = (opt: FlowSequenceStore) => {
    const selectedIds = context.ids ?? [];

    if (selectedIds.length === 0) return;

    if (selectedIds.length === 1) {
      if (contact?.sequence?.id === opt.id || !contact) {
        ui.commandMenu.setOpen(false);

        return;
      }

      opt.linkContact(contact.id, contact.emailId);
    }

    if (selectedIds.length > 1) {
      const selectedContacts = selectedIds
        .map((id) => contacts.value.get(id))
        .filter((contact): contact is ContactStore => contact !== null);

      const hasConflictingSequence = selectedContacts.some(
        (ct) => !!ct.sequence?.id && ct.sequence.id !== opt.id,
      );

      if (hasConflictingSequence) {
        handleOpenConfirmDialog(opt.id);

        return;
      } else {
        flowSequences.linkContacts(opt.id, selectedIds);
      }
    }

    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const sequenceOptions = flowSequences.toComputedArray((arr) => arr);

  const handleCreateOption = (value: string) => {
    flowSequences?.create(
      { name: value, description: '' },
      {
        onSuccess: (sequenceId) => {
          const newSequence = flowSequences.value.get(
            sequenceId,
          ) as FlowSequenceStore;

          if (!newSequence) return;
          handleSelect(newSequence);
        },
      },
    );
  };
  const filteredOptions = sequenceOptions?.filter((sequence) =>
    sequence.value.name.toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <Command shouldFilter={false} label='Move to sequence...'>
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Move to sequence...'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />

      <Command.List>
        {filteredOptions.map((flowSequence) => {
          const isSelected =
            context.ids?.length === 1 &&
            contact?.sequence?.id === flowSequence.id;

          return (
            <CommandItem
              key={flowSequence.id}
              rightAccessory={isSelected ? <Check /> : undefined}
              onSelect={() => {
                handleSelect(flowSequence as FlowSequenceStore);
              }}
            >
              {flowSequence.value.name ?? 'Unnamed'}
            </CommandItem>
          );
        })}

        {search && (
          <CommandItem
            leftAccessory={<Plus />}
            onSelect={() => handleCreateOption(search)}
          >
            <span className='text-gray-700 ml-1'>Create new sequence:</span>
            <span className='text-gray-500 ml-1'>{search}</span>
          </CommandItem>
        )}
      </Command.List>
    </Command>
  );
});
