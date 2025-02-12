import { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { ArchiveAgentUsecase } from '@domain/usecases/agents/archive-agent.usecase';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const ArchiveAgent = observer(() => {
  const { ui } = useStore();
  const navigate = useNavigate();

  const context = ui.commandMenu.context;
  const usecase = useMemo(
    () => new ArchiveAgentUsecase(context.ids[0]),
    [context.ids[0]],
  );

  const navigateToAgentsList = () => {
    navigate('/agents');
  };

  useModKey('Enter', () => {
    return usecase.execute(navigateToAgentsList);
  });

  const handleClose = () => {
    ui.commandMenu.setType('AgentCommands');
    ui.commandMenu.setOpen(false);
  };

  return (
    <Command shouldFilter={false}>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100 cursor-default'>
        <div className='flex justify-between'>
          <h1 className='text-base font-semibold'>
            Archive {usecase.agent?.value.name ?? ''}
          </h1>
          <div>
            <CommandCancelIconButton onClose={handleClose} />
          </div>
        </div>

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='error'
            dataTest={'archive-agent-confirm'}
            onClick={() => usecase.execute(navigateToAgentsList)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                return usecase.execute(navigateToAgentsList);
              }
            }}
          >
            Archive agent
          </Button>
        </div>
      </article>
    </Command>
  );
});
