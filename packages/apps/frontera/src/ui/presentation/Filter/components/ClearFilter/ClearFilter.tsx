import { X } from '@ui/media/icons/X';
import { IconButton } from '@ui/form/IconButton';

interface ClearFilterProps {
  onClearFilter: () => void;
}

export const ClearFilter = ({ onClearFilter }: ClearFilterProps) => {
  return (
    <IconButton
      size='xs'
      variant='outline'
      onClick={onClearFilter}
      colorScheme='grayModern'
      aria-label='clear filter selected'
      icon={<X className='text-gray-500' />}
      className=' bg-white  border-grayModern-300 focus:outline-none'
    />
  );
};
