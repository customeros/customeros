import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { useStore } from '@shared/hooks/useStore';

import { useSlackOauthCallback } from './hooks';
import { Header, AgentCard, EmptyState } from './components';

export const AgentsPage = observer(() => {
  const store = useStore();
  const [firstView, setFirstView] = useLocalStorage(
    'cos_agents_first_view',
    false,
  );

  const agents = store.agents.toArray();

  useSlackOauthCallback();

  if (!firstView) {
    return (
      <EmptyState
        onClick={() => {
          setFirstView(true);
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
            icon={agent.value.icon}
            name={agent.value.name}
            color={agent.value.color}
            status={agent.value.isActive ? 'ON' : 'OFF'}
          />
        ))}
      </div>
    </div>
  );
});
