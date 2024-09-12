import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
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

  return (
    <Command label='Change or add to sequence'>
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Change or add to sequence'
      />

      <Command.List>
        {sequenceOptions.map((flowSequence) => {
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
      </Command.List>
    </Command>
  );
});
