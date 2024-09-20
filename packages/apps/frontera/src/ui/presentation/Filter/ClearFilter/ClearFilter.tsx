import { X } from '@ui/media/icons/X';
import { IconButton } from '@ui/form/IconButton';

export const ClearFilter = () => {
  return (
    <IconButton
      size='xs'
      icon={<X />}
      colorScheme='grayModern'
      className='border-transparent'
      aria-label='clear filter selected'
    />
  );
};
