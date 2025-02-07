import { useParams, useNavigate, useLocation } from 'react-router-dom';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { ButtonGroup } from '@ui/form/ButtonGroup';

export const SubHeader = () => {
  const { id } = useParams<{ id: string }>();

  const location = useLocation();

  const navigate = useNavigate();

  if (!id) {
    throw new Error('No id provided');
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
};
