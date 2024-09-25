import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store.ts';
import { ContactStore } from '@store/Contacts/Contact.store';

import { Plus } from '@ui/media/icons/Plus';
import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditContactFlow = observer(() => {
  const { contacts, ui, flows } = useStore();
  const [search, setSearch] = useState('');

  const context = ui.commandMenu.context;

  const contact = contacts.value.get(context.ids?.[0] as string);
  const selectedIds = context.ids;

  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.value.name}`
      : `${selectedIds?.length} contacts`;

  const handleOpenConfirmDialog = (id: string) => {
    ui.commandMenu.toggle('ConfirmBulkFlowEdit');
    ui.commandMenu.setContext({
      ...ui.commandMenu.context,
      property: id,
    });
    ui.commandMenu.setOpen(true);
  };

  const handleSelect = (opt: FlowStore) => {
    const selectedIds = context.ids ?? [];

    if (selectedIds.length === 0) return;

    if (selectedIds.length === 1) {
      if (contact?.flow?.id === opt.id || !contact) {
        ui.commandMenu.setOpen(false);

        return;
      }

      opt.linkContact(contact.id);
    }

    if (selectedIds.length > 1) {
      const selectedContacts = selectedIds
        .map((id) => contacts.value.get(id))
        .filter((contact): contact is ContactStore => contact !== null);

      const hasConflictingFlow = selectedContacts.some(
        (ct) => !!ct.flow?.id && ct.flow.id !== opt.id,
      );

      if (hasConflictingFlow) {
        handleOpenConfirmDialog(opt.id);

        return;
      } else {
        opt.linkContacts(selectedIds);
      }
    }

    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const flowOptions = flows.toComputedArray((arr) => arr);

  const handleCreateOption = (value: string) => {
    flows?.create(value, {
      onSuccess: (flowId) => {
        const newFlow = flows.value.get(flowId) as FlowStore;

        if (!newFlow) return;
        handleSelect(newFlow);
      },
    });
  };
  const filteredOptions = flowOptions?.filter((flow) =>
    flow.value.name.toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <Command shouldFilter={false} label='Move to flow...'>
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Move to flow...'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />

      <Command.List>
        {filteredOptions.map((flowFlow) => {
          const isSelected =
            context.ids?.length === 1 && contact?.flow?.id === flowFlow.id;

          return (
            <CommandItem
              key={flowFlow.id}
              rightAccessory={isSelected ? <Check /> : undefined}
              onSelect={() => {
                handleSelect(flowFlow as FlowStore);
              }}
            >
              {flowFlow.value.name ?? 'Unnamed'}
            </CommandItem>
          );
        })}

        {search && (
          <CommandItem
            leftAccessory={<Plus />}
            onSelect={() => handleCreateOption(search)}
          >
            <span className='text-gray-700 ml-1'>Create new flow:</span>
            <span className='text-gray-500 ml-1'>{search}</span>
          </CommandItem>
        )}
      </Command.List>
    </Command>
  );
});
