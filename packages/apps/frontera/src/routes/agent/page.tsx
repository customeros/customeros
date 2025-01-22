import { useMemo, useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { AgentViewUsecase } from '@domain/usecases/agents/agent-view.usecase';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { Settings02 } from '@ui/media/icons/Settings02';

import { Header, capabilities } from './components';

const usecase = new AgentViewUsecase();

export const AgentPage = observer(() => {
  const store = useStore();
  const { id } = useParams<{ id: string }>();

  const agent = id ? store.agents.getById(id) : null;

  const ActiveCapability = useMemo(
    () => capabilities[usecase.activeCapability.type],
    [usecase.activeCapability.type],
  );

  useEffect(() => {
    if (agent) {
      usecase.setActiveCapability(agent.value.capabilities[0]);
    }
  }, []);

  if (!id) {
    throw new Error('No id provided');
  }

  return (
    <div>
      <Header />

      <div className='flex h-screen'>
        <div className='w-[448px] border-r border-r-grayModern-200 px-4 py-3'>
          <div className='mb-2'>
            <h2 className='font-medium mb-1'>Goal</h2>
            <p className='pb-2 text-sm'>{agent?.value.goal}</p>
          </div>

          <h2 className='mb-2 font-medium text-sm'>This agent will:</h2>

          <ul className='space-y-1'>
            {agent?.value.capabilities.map((capability) => (
              <div
                key={capability.id}
                onClick={() => {
                  if (capability.values.length) {
                    usecase.setActiveCapability(capability);
                  }
                }}
                className={cn(
                  'flex items-center px-2 py-1 justify-between rounded-lg select-none',
                  capability.values.length &&
                    'hover:bg-grayModern-100 cursor-pointer',
                  capability.id === usecase.activeCapability.id &&
                    'bg-grayModern-200 hover:bg-grayModern-200',
                )}
              >
                <div className='flex items-center gap-2'>
                  <DotSingle className='text-grayModern-500' />
                  <p>{capability.name}</p>
                </div>

                {capability.values.length && (
                  <Settings02 className='text-grayModern-500' />
                )}
              </div>
            ))}
          </ul>
        </div>

        {usecase.activeCapability && (
          <div className='w-[418px] border-r border-r-grayModern-200 px-4 py-3'>
            <ActiveCapability />
          </div>
        )}
      </div>
    </div>
  );
});
