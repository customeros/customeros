import { useMemo } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { RenameAgentUsecase } from '@domain/usecases/agents/rename-agent.usecase';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const RenameAgent = observer(() => {
  const { ui } = useStore();
  const context = ui.commandMenu.context;
  const usecase = useMemo(
    () => new RenameAgentUsecase(context.ids[0]),
    [context.ids[0]],
  );

  useKey('Escape', () => {
    ui.commandMenu.setType('AgentCommands');
    ui.commandMenu.setOpen(false);
  });

  return (
    <Command shouldFilter={false} label={`Rename agent`}>
      <CommandInput
        value={usecase.inputValue}
        placeholder='Rename agent'
        label={`Agent - ${usecase.agent?.value.name}`}
        onValueChange={(value) => usecase.setInputValue(value)}
      />
      <Command.List>
        <CommandItem
          onSelect={usecase.execute}
          leftAccessory={<Edit03 />}
        >{`Rename agent to "${usecase.inputValue}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
