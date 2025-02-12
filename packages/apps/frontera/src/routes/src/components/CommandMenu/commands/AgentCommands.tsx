import { observer } from 'mobx-react-lite';

import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const AgentCommands = observer(() => {
  const store = useStore();
  const id = store.ui.commandMenu.context.ids?.[0];
  const agent = id ? store.agents.getById(id) : null;
  const label = `Agent - ${agent?.value.name}`;

  return (
    <CommandsContainer label={label}>
      <CommandItem
        leftAccessory={<Icon name='edit-03' />}
        onSelect={() => {
          store.ui.commandMenu.setType('RenameAgent');
        }}
      >
        Rename agent
      </CommandItem>
      {/*<CommandItem*/}
      {/*  leftAccessory={<Icon name='layers-two-01' />}*/}
      {/*  onSelect={() => {*/}
      {/*    store.ui.commandMenu.setType('DuplicateAgent');*/}
      {/*  }}*/}
      {/*>*/}
      {/*  Duplicate agent*/}
      {/*</CommandItem>*/}

      <CommandItem
        leftAccessory={<Icon name='layers-two-01' />}
        onSelect={() => {
          store.ui.commandMenu.setType('ArchiveAgent');
        }}
      >
        Archive agent
      </CommandItem>
    </CommandsContainer>
  );
});
