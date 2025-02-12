import { useRef, useMemo, useEffect } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { DuplicateAgentUsecase } from '@domain/usecases/agents/duplicate-agent.usecase';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const DuplicateAgent = observer(() => {
  const { ui } = useStore();
  const inputRef = useRef<HTMLInputElement>(null);

  const context = ui.commandMenu.context;
  const usecase = useMemo(
    () => new DuplicateAgentUsecase(context.ids[0]),
    [context.ids[0]],
  );

  useEffect(() => {
    if (inputRef.current) {
      setTimeout(() => {
        inputRef.current?.focus();
        inputRef.current?.select();
      }, 50);
    }
  }, []);

  useKey('Escape', () => {
    ui.commandMenu.setType('AgentCommands');
    ui.commandMenu.setOpen(false);
  });

  if (!usecase.agent) {
    return null;
  }

  return (
    <Command shouldFilter={false} label={`Duplicate agent`}>
      <div className='p-6 pb-4 flex flex-col gap-1 '>
        <div className='flex items-center justify-between'>
          <h1
            className='text-base font-semibold'
            data-test='create-new-flow-modal-title'
          >
            Duplicate {usecase.agent.value.name ?? ''}
          </h1>

          <CommandCancelIconButton
            onClose={() => {
              ui.commandMenu.setOpen(false);
              ui.commandMenu.clearContext();
            }}
          />
        </div>
      </div>

      <div className='pr-6 pl-6 pb-6 flex flex-col gap-2 '>
        <Input
          autoFocus
          size='sm'
          ref={inputRef}
          id='agent-name'
          variant='unstyled'
          placeholder='Agent name'
          value={usecase.inputValue}
          dataTest='duplicated-agent-name'
          onChange={(e) => {
            usecase.setInputValue(e.target.value);
          }}
          onKeyDown={(e) => {
            e.stopPropagation();

            if (e.key === 'Enter') {
              usecase.execute();
            }

            if (e.key === 'Escape') {
              ui.commandMenu.setType('AgentCommands');
              ui.commandMenu.setOpen(false);
            }
          }}
        />
      </div>

      <div className='flex w-full gap-2 pl-6 pr-6 pb-6'>
        <CommandCancelButton
          dataTest='cancel-create-new-flow'
          onClose={() => {
            ui.commandMenu.setOpen(false);
            ui.commandMenu.clearContext();
          }}
        />

        <Button
          className='w-full'
          colorScheme='primary'
          isLoading={usecase.isSaving}
          onClick={() => usecase.execute()}
          loadingText={'Duplicating agent…'}
          data-test='confirm-duplicate-agent'
        >
          Duplicate Agent
        </Button>
      </div>
    </Command>
  );
});
