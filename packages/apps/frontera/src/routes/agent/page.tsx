import { useMemo, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import get from 'lodash/get';
import { observer } from 'mobx-react-lite';
import { AgentViewUsecase } from '@domain/usecases/agents/agent-view.usecase';

import { cn } from '@ui/utils/cn';
import { Icon, IconName } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';

import { Header, capabilities } from './components';

const goalMap = {
  identify_web_visitor:
    'Identify website visitors and add them as enriched leads to CustomerOS',
  evaluate_icp_fit:
    'Qualify new leads based on whether they match your ideal customer profile or not',
};

export const AgentPage = observer(() => {
  const store = useStore();

  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [queryParams, setQueryParams] = useSearchParams();

  const agent = id ? store.agents.getById(id) : null;

  const usecase = useMemo(
    () => new AgentViewUsecase(id ?? '', queryParams.get('cid')),
    [id],
  );

  const ActiveCapability = useMemo(
    () =>
      usecase.activeCapability
        ? capabilities[usecase.activeCapability.type]
        : () => null,
    [usecase?.activeCapability?.type],
  );

  useEffect(() => {
    if (!queryParams.get('cid')) {
      setQueryParams((params) => {
        if (!usecase.activeCapability) return params;
        params.set('cid', usecase.activeCapability?.id);

        return params;
      });
    }
  }, []);

  if (!id) {
    throw new Error('No id provided');
  }

  return (
    <div>
      <Header
        isActive={!!agent?.value.isActive}
        agentName={agent?.value?.name ?? ''}
        icon={agent?.value.icon as IconName}
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
                    navigate(`?cid=${capability.id}`);
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
                  <p className='text-sm'>{capability.name ?? 'Unknown'}</p>
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
