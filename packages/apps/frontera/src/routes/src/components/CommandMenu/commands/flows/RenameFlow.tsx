import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const RenameFlow = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const entity = store.flows.value.get(context.ids?.[0] as string);
  const label = `Flow - ${entity?.value?.name}`;
  const defaultValue = entity?.value?.name ?? '';

  const [name, setName] = useState(() => defaultValue ?? '');

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('FlowCommands');
  };

  const handleSelect = () => {
    if (!context.ids?.[0]) return;

    if (!entity) return;

    if (!name.trim()?.length) {
      handleClose();

      return;
    }

    entity?.update((value) => {
      value.name = name;

      return value;
    });
    handleClose();
  };

  return (
    <Command label={`Rename flow`}>
      <CommandInput
        value={name}
        label={label}
        placeholder='Rename flow'
        onValueChange={(value) => setName(value)}
      />
      <Command.List>
        <CommandItem
          onSelect={handleSelect}
          leftAccessory={<Edit03 />}
        >{`Rename flow to "${name}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
