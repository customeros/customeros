import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { TextInput } from '@ui/media/icons/TextInput.tsx';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const RenameTableViewDef = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset || '');
  const tableViewName = tableViewDef?.value.name;
  const [value, setValue] = useState(() => tableViewName ?? '');

  const label = `View - ${tableViewName}`;

  const handleSelect = () => {
    tableViewDef?.update((opp) => {
      opp.name = value;

      return opp;
    });
    store.ui.commandMenu.setOpen(false);
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
