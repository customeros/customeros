import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';

export const PersonasFlowProfileMenu = () => {
  return (
    <div className='flex justify-between'>
      <p className='font-semibold'>Flows</p>
      <IconButton
        size='xs'
        icon={<Plus />}
        variant='ghost'
        aria-label='add new flow'
      />
    </div>
  );
};
