import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store';

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

  const handleSelect = (opt: FlowSequenceStore) => {
    if (!context.ids?.[0] || !contact) return;

    if (selectedIds?.length === 1) {
      opt.linkContact(contact.id, contact.emailId);
    }

    if (selectedIds?.length > 1) {
      flowSequences.linkContacts(opt.id, selectedIds);
    }

    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  return (
    <Command label='Change or add to sequence'>
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Change or add to sequence'
      />

      <Command.List>
        {flowSequences.toArray().map((flowSequence) => {
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
              {flowSequence.value.name}
            </CommandItem>
          );
        })}
      </Command.List>
    </Command>
  );
});
