import { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';

import { CreateAgentUsecase } from '@domain/usecases/agents/create-agent.usecase';

import { Icon } from '@ui/media/Icon';
import { AgentType } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const CreateAgent = () => {
  const navigate = useNavigate();

  const navigateToAgent = (id: string) => {
    navigate(`/agents/${id}`);
  };
  const usecase = useMemo(() => new CreateAgentUsecase(navigateToAgent), []);

  return (
    <Command>
      <CommandInput
        label='Add a new agent'
        placeholder='Search for an agent...'
      />
      <Command.List>
        <CommandItem
          onSelect={() => usecase.execute(AgentType.WebVisitIdentifier)}
        >
          <Icon name='radar' />
          <span>Web visitor identifier</span>
        </CommandItem>
        <CommandItem onSelect={() => usecase.execute(AgentType.IcpQualifier)}>
          <Icon name='target-04' />
          <span>ICP qualifier</span>
        </CommandItem>
        <CommandItem
          onSelect={() => usecase.execute(AgentType.CampaignManager)}
        >
          <Icon name='send-03' />
          <span>Outbound campaign manager</span>
        </CommandItem>
        <CommandItem onSelect={() => usecase.execute(AgentType.MeetingKeeper)}>
          <Icon name='edit-04' />
          <span>Meeting keeper</span>
        </CommandItem>
        <CommandItem onSelect={() => usecase.execute(AgentType.SupportSpotter)}>
          <Icon name='life-buoy-01' />
          <span>Support spotter</span>
        </CommandItem>
      </Command.List>
    </Command>
  );
};
