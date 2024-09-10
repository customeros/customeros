import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const RenameOpportunityName = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const opportunity = store.opportunities.value.get(context.ids?.[0]);
  const [value, setValue] = useState(() => opportunity?.value.name ?? '');

  const label = `Opportunity - ${opportunity?.value.name}`;

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('OpportunityCommands');
  };

  const handleSelect = () => {
    if (!value.trim()?.length) {
      handleClose();

      return;
    }

    opportunity?.update((opp) => {
      opp.name = value;

      return opp;
    });
  };

  return (
    <Command shouldFilter={false}>
      <CommandInput
        value={value}
        label={label}
        onValueChange={setValue}
        placeholder='Rename opportunity'
      />
      <Command.List>
        <CommandItem
          onSelect={handleSelect}
          leftAccessory={<Edit03 />}
        >{`Rename opportunity to "${value}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
