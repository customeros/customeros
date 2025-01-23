import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

import { Header, AgentCard, EmptyState } from './components';

export const AgentsPage = observer(() => {
  const store = useStore();
  const [isVisible, setIsVisible] = useState(true);

  const agents = store.agents.toArray();

  if (isVisible) {
    return (
      <EmptyState
        onClick={() => {
          setIsVisible(!isVisible);
        }}
      />
    );
  }

  return (
    <div className='relative h-full'>
      <Header />
      <div className='columns-3 p-4 gap-4'>
        {agents.map((agent) => (
          <AgentCard
            id={agent.id}
            key={agent.id}
            name={agent.value.name}
            status={agent.value.isActive ? 'ON' : 'OFF'}
          />
        ))}
      </div>
    </div>
  );
});
