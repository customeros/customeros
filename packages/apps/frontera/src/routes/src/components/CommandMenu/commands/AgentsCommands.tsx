import { observer } from 'mobx-react-lite';

import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const AgentsCommands = observer(() => {
  const store = useStore();

  return (
    <CommandsContainer label={'Agents'}>
      <CommandItem
        leftAccessory={<Icon name='plus-circle' />}
        onSelect={() => {
          store.ui.commandMenu.setType('CreateAgent');
        }}
      >
        Add new agent...
      </CommandItem>
    </CommandsContainer>
  );
});
