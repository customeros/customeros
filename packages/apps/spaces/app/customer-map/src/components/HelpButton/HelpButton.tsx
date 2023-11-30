import { IconButton } from '@ui/form/IconButton';
import { HelpCircle } from '@ui/media/icons/HelpCircle';

interface HelpButtonProps {
  isOpen: boolean;
  onOpen: () => void;
}

export const HelpButton = ({ isOpen, onOpen }: HelpButtonProps) => {
  return (
    <IconButton
      size='xs'
      variant='ghost'
      onClick={onOpen}
      aria-label='Help'
      id='help-button'
      icon={<HelpCircle color='gray.400' />}
      visibility={isOpen ? 'visible' : 'hidden'}
    />
  );
};
