import { observer } from 'mobx-react-lite';

import { Icon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';

export const Header = observer(() => {
  const store = useStore();

  return (
    <div className='w-full border-b border-b-gray-200 px-3 py-[8px] flex justify-between items-center'>
      <p className='text-md font-medium'>Agents</p>

      <Button
        size='xs'
        variant='outline'
        colorScheme='primary'
        onClick={() => {
          store.ui.commandMenu.toggle('CreateAgent');
        }}
      >
        <Icon name='plus' />
        Add agent
      </Button>
    </div>
  );
});
