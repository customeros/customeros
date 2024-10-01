import { observer } from 'mobx-react-lite';

import { Plus } from '@ui/media/icons/Plus.tsx';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';

export const CreateSequenceButton = observer(() => {
  const store = useStore();

  const handleCreateSequence = () => {
    store.ui.commandMenu.toggle('CreateNewFlow');
  };

  return (
    <Button
      size='xs'
      className='mr-1'
      leftIcon={<Plus />}
      colorScheme='primary'
      dataTest='add-new-flow'
      onClick={handleCreateSequence}
    >
      New flow
    </Button>
  );
});
