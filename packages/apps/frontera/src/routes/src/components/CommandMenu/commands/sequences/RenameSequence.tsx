import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const RenameSequence = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const entity = store.flowSequences.value.get(context.ids?.[0] as string);
  const label = `Sequence - ${entity?.value?.name}`;

  const handleSelect = () => {
    if (!context.ids?.[0]) return;

    if (!entity) return;

    entity?.update((value) => {
      value.name = name;

      return value;
    });
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('SequenceCommands');
  };

  const defaultValue = entity?.value?.name ?? '';

  const [name, setName] = useState(() => defaultValue ?? '');

  return (
    <Command
      label={`Rename sequence`}
      onKeyDown={(e) => {
        e.stopPropagation();
      }}
    >
      <CommandInput
        value={name}
        label={label}
        placeholder='Rename sequence...'
        onValueChange={(value) => setName(value)}
      />
      <Command.List>
        <CommandItem
          onSelect={handleSelect}
          leftAccessory={<Edit03 />}
        >{`Rename sequence to "${name}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
