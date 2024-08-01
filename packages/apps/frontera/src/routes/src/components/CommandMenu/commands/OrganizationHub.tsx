import { observer } from 'mobx-react-lite';

import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

import { GlobalSharedCommands } from './GlobalHub';

export const OrganizationHub = observer(() => {
  return (
    <Command>
      <CommandInput
        label='Organizations'
        placeholder='Type a command or search'
      />
      <Command.List>
        <CommandItem onSelect={() => {}} leftAccessory={<PlusCircle />}>
          Add new organizations...
        </CommandItem>

        <Command.Group heading='Navigate'>
          <GlobalSharedCommands />
        </Command.Group>
      </Command.List>
    </Command>
  );
});
