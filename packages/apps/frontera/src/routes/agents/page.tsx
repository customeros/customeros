import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { useStore } from '@shared/hooks/useStore';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';

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

  useEffect(() => {
    store.ui.commandMenu.setType('AgentsCommands');
  }, []);

  return (
    <div className='relative h-full'>
      <Header />
      <ScrollAreaRoot className='h-full'>
        <ScrollAreaViewport>
          <div className='flex flex-wrap p-4 gap-4'>
            {agents
              .filter((agent) => agent.value.visible)
              .map((agent) => (
                <AgentCard
                  id={agent.id}
                  key={agent.id}
                  icon={agent.value.icon}
                  name={agent.value.name}
                  colorMap={agent.colorMap}
                  defaultName={agent.defaultName}
                  status={agent.value.isActive ? 'ON' : 'OFF'}
                  hasError={!!agent.value.error || !agent.value.isConfigured}
                />
              ))}
            <div className='min-w-[372px] flex-1 p-3'></div>
            <div className='min-w-[372px] flex-1 p-3'></div>
            <div className='min-w-[372px] flex-1 p-3'></div>
          </div>
        </ScrollAreaViewport>
        <ScrollAreaScrollbar orientation='vertical'>
          <ScrollAreaThumb />
        </ScrollAreaScrollbar>
      </ScrollAreaRoot>
    </div>
  );
});
