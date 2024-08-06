import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

import { GlobalSharedCommands } from './GlobalHub';

export const OrganizationHub = observer(() => {
  const store = useStore();

  return (
    <Command>
      <CommandInput
        label='Organizations'
        placeholder='Type a command or search'
      />
      <Command.List>
        <CommandItem
          leftAccessory={<PlusCircle />}
          onSelect={() => {
            store.ui.commandMenu.setType('AddNewOrganization');
          }}
        >
          Add new organizations...
        </CommandItem>

        <Command.Group heading='Navigate'>
          <GlobalSharedCommands />
        </Command.Group>
      </Command.List>
    </Command>
  );
});
