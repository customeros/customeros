import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';

export const PersonasFlowProfileMenu = () => {
  return (
    <div className='flex justify-between'>
      <p className='font-semibold'>Flows</p>
      <IconButton
        icon={<Plus />}
        aria-label='add new flow'
        variant='ghost'
        size='xs'
      />
    </div>
  );
};
