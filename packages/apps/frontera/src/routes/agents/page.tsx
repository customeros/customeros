import { useState } from 'react';

import { Header, AgentCard, EmptyState } from './components';

export const AgentsPage = () => {
  const [isVisible, setIsVisible] = useState(true);

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
        <AgentCard status='ON' name='Web visit identifier' />
        <AgentCard status='OFF' name='Lead qualifier' />
        <AgentCard status='OFF' name='Manager: Cold outbound example' />
        <AgentCard status='OFF' name='Lead qualifier' />
        <AgentCard status='ON' name='Web visit identifier' />
      </div>
    </div>
  );
};
