import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { AgentViewUsecase } from '@domain/usecases/agents/agent-view.usecase';

import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const AgentCommands = observer(() => {
  const store = useStore();
  const id = store.ui.commandMenu.context.ids?.[0];
  const agent = id ? store.agents.getById(id) : null;
  const label = `Agent - ${agent?.value.name}`;
  const usecase = useMemo(() => new AgentViewUsecase(id ?? ''), [id]);

  return (
    <CommandsContainer label={label}>
      <CommandItem
        leftAccessory={<Icon name='edit-03' />}
        onSelect={() => {
          store.ui.commandMenu.setType('RenameAgent');
        }}
        rightAccessory={
          <>
            <Kbd>
              <Icon className={'size-3'} name='arrow-block-up' />
            </Kbd>
            <Kbd>R</Kbd>
          </>
        }
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
        leftAccessory={<Icon name='columns-03' />}
        onSelect={() => {
          usecase.toggleActive();
          store.ui.commandMenu.setOpen(false);
        }}
        rightAccessory={
          <>
            <Kbd>
              <Icon className={'size-3'} name='arrow-block-up' />
            </Kbd>
            <Kbd>S</Kbd>
          </>
        }
      >
        Change agent status
      </CommandItem>

      <CommandItem
        leftAccessory={<Icon name='layers-two-01' />}
        onSelect={() => {
          store.ui.commandMenu.setType('ArchiveAgent');
        }}
        rightAccessory={
          <>
            <CommandKbd />
            <Kbd>
              <Icon name='delete' className='size-3' />
            </Kbd>
          </>
        }
      >
        Archive agent
      </CommandItem>
    </CommandsContainer>
  );
});
