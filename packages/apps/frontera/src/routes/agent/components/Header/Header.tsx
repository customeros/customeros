import { useMemo } from 'react';
import { Link, useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { ToggleAgentActiveUsecase } from '@domain/usecases/agents/toggle-agent-active.usecase';

import { Switch } from '@ui/form/Switch';
import { Icon, IconName } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';

export const Header = observer(() => {
  const store = useStore();

  const { id } = useParams<{ id: string }>();

  const agent = id ? store.agents.getById(id) : null;

  const usecase = useMemo(() => new ToggleAgentActiveUsecase(id ?? ''), [id]);

  if (!id) {
    throw new Error('No id provided');
  }

  return (
    <div className='w-full border-b border-b-gray-200 px-3 py-[11px] flex justify-between items-center'>
      <div className='flex items-center gap-1'>
        <Link
          to='/agents'
          className='text-md font-medium text-grayModern-500 hover:text-grayModern-700 transition-colors'
        >
          Agents
        </Link>
        <Icon name='chevron-right' className='w-4 h-4 text-gray-500' />
        <div className='flex items-center gap-1'>
          {agent?.value?.icon && (
            <div className='flex items-center rounded-sm p-0.5 bg-grayModern-100'>
              <Icon
                name={agent.value.icon as IconName}
                className='w-4 h-4 text-grayModern-500'
              />
            </div>
          )}

          <div className='flex items-center gap-1'>
            <p className='text-md font-medium'>{agent?.value?.name}</p>
          </div>

          <div className='ml-4 flex items-center'>
            <Switch
              checked={agent?.value?.isActive}
              onChange={() => usecase.toggleActive()}
            />
          </div>
        </div>
      </div>
    </div>
  );
});
