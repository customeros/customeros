import { useParams, useLocation, useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { AgentType } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';

export const SubHeader = observer(() => {
  const { id } = useParams<{ id: string }>();

  const location = useLocation();
  const { agents } = useStore();
  const navigate = useNavigate();

  const agent = id ? agents.getById(id) : null;

  if (!id) {
    throw new Error('No id provided');
  }

  if (agent?.type !== AgentType.CampaignManager) {
    return null;
  }

  return (
    <div className='border py-2 px-2'>
      <ButtonGroup className='flex items-center w-[303px] '>
        <Button
          size='xs'
          onClick={() => navigate(`/agents/${id}/setup`)}
          leftIcon={<Icon name='settings-02' className='text-inherit' />}
          className={cn('w-full', {
            selected: location.pathname.includes('/setup'),
          })}
        >
          Setup
        </Button>
        <Button
          size='xs'
          onClick={() => navigate(`/agents/${id}/editor`)}
          leftIcon={<Icon name='arrow-if-path' className='text-inherit' />}
          className={cn('w-full', {
            selected: location.pathname.includes('/editor'),
          })}
        >
          Editor
        </Button>
        <Button
          size='xs'
          onClick={() => navigate(`/agents/${id}/list`)}
          leftIcon={<Icon name='users-02' className='text-inherit' />}
          className={cn('w-full', {
            selected: location.pathname.includes('/list'),
          })}
        >
          List
        </Button>
      </ButtonGroup>
    </div>
  );
});
