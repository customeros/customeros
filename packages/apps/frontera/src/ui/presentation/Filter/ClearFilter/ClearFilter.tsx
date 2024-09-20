import { X } from '@ui/media/icons/X';
import { IconButton } from '@ui/form/IconButton';

interface ClearFilterProps {
  onClearFilter: () => void;
}

export const ClearFilter = ({ onClearFilter }: ClearFilterProps) => {
  return (
    <IconButton
      size='xs'
      icon={<X />}
      onClick={onClearFilter}
      colorScheme='grayModern'
      className='border-transparent'
      aria-label='clear filter selected'
    />
  );
};
