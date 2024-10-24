import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { Command, CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

import { GlobalSharedCommands } from './GlobalHub';

export const OrganizationHub = observer(() => {
  const store = useStore();

  return (
    <CommandsContainer label={'Organizations'} dataTest={'organization-hub'}>
      <CommandItem
        leftAccessory={<PlusCircle />}
        dataTest={'organization-hub-add-new-orgs'}
        onSelect={() => {
          store.ui.commandMenu.setType('AddNewOrganization');
        }}
      >
        Add new organizations...
      </CommandItem>
      <Command.Group heading='Navigate'>
        <GlobalSharedCommands />
      </Command.Group>
    </CommandsContainer>
  );
});
