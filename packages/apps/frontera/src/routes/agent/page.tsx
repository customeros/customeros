import { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import get from 'lodash/get';
import { observer } from 'mobx-react-lite';
import { AgentViewUsecase } from '@domain/usecases/agents/agent-view.usecase';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';

import { Header, capabilities } from './components';

const goalMap = {
  identify_web_visitor:
    'Identify website visitors and add them as enriched leads to CustomerOS',
};

export const AgentPage = observer(() => {
  const store = useStore();
  const { id } = useParams<{ id: string }>();

  const agent = id ? store.agents.getById(id) : null;

  const usecase = useMemo(() => new AgentViewUsecase(id ?? ''), [id]);

  const ActiveCapability = useMemo(
    () =>
      usecase.activeCapability
        ? capabilities[usecase.activeCapability.type]
        : () => null,
    [usecase?.activeCapability?.type],
  );

  if (!id) {
    throw new Error('No id provided');
  }

  return (
    <div>
      <Header
        isActive={!!agent?.value.isActive}
        agentName={agent?.value?.name ?? ''}
        onToggleActive={usecase.toggleActive}
      />

      <div className='flex h-screen'>
        <div className='w-[448px] border-r border-r-grayModern-200 px-4 py-3'>
          <div className='mb-2'>
            <h2 className='font-medium mb-1'>Goal</h2>
            <p className='pb-2 text-sm'>
              {get(goalMap, agent?.value.goal ?? '', 'Unknown')}
            </p>
          </div>

          <h2 className='mb-2 font-medium text-sm'>This agent will:</h2>

          <ul className='space-y-1'>
            {agent?.value.capabilities.map((capability) => (
              <div
                key={capability.id}
                onClick={() => {
                  if (capability.config.length) {
                    usecase.setActiveCapability(capability);
                  }
                }}
                className={cn(
                  'flex items-center px-2 py-1 justify-between rounded-lg select-none',
                  capability.config.length &&
                    'hover:bg-grayModern-200 cursor-pointer',
                  capability.id === usecase?.activeCapability?.id &&
                    'bg-grayModern-100 hover:bg-grayModern-100',
                )}
              >
                <div className='flex items-center gap-2'>
                  <Icon
                    stroke={capability.errors ? 'currentColor' : 'none'}
                    name={capability.errors ? 'radio-dot' : 'dot-single'}
                    className={
                      capability.errors
                        ? 'text-error-500'
                        : 'text-grayModern-500'
                    }
                  />
                  <p>{capability.name ?? 'Unknown'}</p>
                </div>

                {capability.config.length > 0 && (
                  <Icon name='settings-02' className='text-grayModern-500' />
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
