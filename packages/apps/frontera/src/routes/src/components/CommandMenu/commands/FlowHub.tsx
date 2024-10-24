import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { Command, CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

import { GlobalSharedCommands } from './GlobalHub';

export const FlowHub = observer(() => {
  const store = useStore();

  return (
    <CommandsContainer label={'Flows'}>
      <CommandItem
        leftAccessory={<PlusCircle />}
        onSelect={() => {
          store.ui.commandMenu.setType('CreateNewFlow');
        }}
      >
        Add new flow...
      </CommandItem>
      <Command.Group heading='Navigate'>
        <GlobalSharedCommands />
      </Command.Group>
    </CommandsContainer>
  );
});
