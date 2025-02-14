import { useMemo, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { useKeys } from 'rooks';
import { observer } from 'mobx-react-lite';
import { AgentViewUsecase } from '@domain/usecases/agents/agent-view.usecase';

import { cn } from '@ui/utils/cn';
import { Icon, IconName } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';

import { Scope } from './components/Scope';
import { goals, Header, configs } from './components';

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

  const ActiveConfig = useMemo(
    () =>
      usecase.activeConfig ? configs[usecase.activeConfig.type] : () => null,
    [usecase?.activeConfig?.type],
  );

  useEffect(() => {
    if (!queryParams.get('cid')) {
      setQueryParams(
        (params) => {
          if (!usecase.activeConfig?.id) return params;
          params.set('cid', usecase.activeConfig.id);

          return params;
        },
        { replace: true },
      );
    }
  }, []);

  useEffect(() => {
    if (id) {
      store.ui.commandMenu.setContext({
        ...store.ui.commandMenu.context,
        ids: [id],
      });
      store.ui.commandMenu.setType('AgentCommands');
    }
  }, [id]);

  useKeys(['Shift', 'S'], (e) => {
    e.stopPropagation();
    e.preventDefault();
    usecase.toggleActive();
  });

  useKeys(['Shift', 'R'], (e) => {
    e.stopPropagation();
    e.preventDefault();
    store.ui.commandMenu.setType('RenameAgent');
    store.ui.commandMenu.setOpen(true);
  });

  useModKey('Backspace', () => {
    store.ui.commandMenu.setType('ArchiveAgent');
    store.ui.commandMenu.setOpen(true);
  });

  if (!id) {
    throw new Error('No id provided');
  }

  if (!agent) {
    return null;
  }

  return (
    <div>
      <Header
        id={id}
        colorMap={agent.colorMap}
        isActive={!!agent?.value.isActive}
        agentName={agent?.value?.name ?? ''}
        icon={agent?.value.icon as IconName}
        onToggleActive={usecase.toggleActive}
      />

      <div className='flex h-screen'>
        <div className='w-[448px] border-r border-r-grayModern-200 px-4 py-3'>
          <div className='mb-2'>
            <div className='mb-1 flex items-center justify-between'>
              <h2 className='font-medium text-sm'>About this agent</h2>
              {agent?.value.scope && <Scope scope={agent.value.scope} />}
            </div>
            <p className='pb-2 text-sm'>{agent?.value.goal ?? 'Unknown'}</p>
          </div>

          <h2 className='mb-2 font-medium text-sm'>It's goal is to</h2>
          <ul className='space-y-1 mb-4'>
            {agent &&
              goals[agent?.value.type]?.map((goal) => (
                <div
                  key={goal}
                  className='flex items-center px-2 py-1 justify-between rounded-lg select-none'
                >
                  <div className='flex items-center gap-2'>
                    <Icon
                      stroke='none'
                      name='dot-single'
                      className={'text-grayModern-500'}
                    />
                    <p className='text-sm'>{goal}</p>
                  </div>
                </div>
              ))}
          </ul>

          <h2 className='mb-2 font-medium text-sm'>It listens for</h2>
          <ul className='space-y-1 mb-4'>
            {agent?.value.listeners.map((listener) => (
              <div
                key={listener.id}
                onClick={() => {
                  if (listener.config.length) {
                    usecase.setActiveConfig(listener);
                    navigate(`?lid=${listener.id}`, { replace: true });
                  }
                }}
                className={cn(
                  'flex items-center px-2 py-1 justify-between rounded-lg select-none',
                  listener.config.length &&
                    'hover:bg-grayModern-200 cursor-pointer',
                  listener.id === usecase?.activeConfig?.id &&
                    'bg-grayModern-100 hover:bg-grayModern-100 font-medium',
                )}
              >
                <div className='flex items-center gap-2'>
                  <Icon
                    stroke={listener.errors ? 'currentColor' : 'none'}
                    name={listener.errors ? 'radio-dot' : 'dot-single'}
                    className={
                      listener.errors ? 'text-error-500' : 'text-grayModern-500'
                    }
                  />
                  <p className='text-sm'>{listener.name ?? 'Unknown'}</p>
                </div>

                {listener.config.length > 0 && (
                  <Icon name='settings-02' className='text-grayModern-500' />
                )}
              </div>
            ))}
          </ul>

          <h2 className='mb-2 font-medium text-sm'>
            It can perform these actions
          </h2>

          <ul className='space-y-1'>
            {agent?.value.capabilities.map((capability) => (
              <div
                key={capability.id}
                onClick={() => {
                  if (capability.config.length) {
                    usecase.setActiveConfig(capability);
                    navigate(`?cid=${capability.id}`, { replace: true });
                  }
                }}
                className={cn(
                  'flex items-center px-2 py-1 justify-between rounded-lg select-none',
                  capability.config.length &&
                    'hover:bg-grayModern-200 cursor-pointer',
                  capability.id === usecase?.activeConfig?.id &&
                    'bg-grayModern-100 hover:bg-grayModern-100 font-medium',
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

        {usecase.activeConfig && ActiveConfig && (
          <div className='w-[418px] border-r border-r-grayModern-200 px-4 py-3'>
            <ActiveConfig />
          </div>
        )}
      </div>
    </div>
  );
});
