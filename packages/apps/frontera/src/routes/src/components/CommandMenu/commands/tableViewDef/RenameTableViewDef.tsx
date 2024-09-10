import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { TextInput } from '@ui/media/icons/TextInput.tsx';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const RenameTableViewDef = observer(() => {
  const store = useStore();

  const id = store.ui.commandMenu.context.ids?.[0] ?? '';
  const tableViewDef = store.tableViewDefs.getById(id);
  const tableViewName = tableViewDef?.value.name;
  const [value, setValue] = useState(() => tableViewName ?? '');

  const label = `View - ${tableViewName}`;

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
  };

  const handleSelect = () => {
    if (!value.trim()?.length) {
      handleClose();

      return;
    }

    tableViewDef?.update((opp) => {
      opp.name = value;

      return opp;
    });
    handleClose();
  };

  return (
    <Command shouldFilter={false}>
      <CommandInput
        value={value}
        label={label}
        onValueChange={setValue}
        placeholder='Rename table view'
      />
      <Command.List>
        <CommandItem
          onSelect={handleSelect}
          leftAccessory={<TextInput />}
        >{`Rename "${value}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
