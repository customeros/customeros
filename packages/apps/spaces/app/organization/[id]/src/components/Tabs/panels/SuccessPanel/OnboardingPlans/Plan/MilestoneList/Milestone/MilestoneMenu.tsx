import { useDisclosure } from '@ui/utils';
import { IconButton } from '@ui/form/IconButton';
import { MinusCircle } from '@ui/media/icons/MinusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';

interface MilestoneMenuProps {
  opacity?: number;
  transition?: string;
  isOptional?: boolean;
  onRetire?: () => void;
  onDuplicate?: () => void;
  onMakeOptional?: () => void;
}

export const MilestoneMenu = ({
  onRetire,
  isOptional,
  onDuplicate,
  onMakeOptional,
  ...buttonProps
}: MilestoneMenuProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Menu isOpen={isOpen} onClose={onClose} onOpen={onOpen}>
      <MenuButton
        as={IconButton}
        size='xs'
        variant='ghost'
        aria-label='Add milestones'
        icon={<DotsVertical color='gray.400' />}
        {...buttonProps}
        opacity={isOpen ? 1 : buttonProps.opacity}
      />
      <MenuList minW='10rem'>
        <MenuItem onClick={onRetire} icon={<MinusCircle color='gray.500' />}>
          Remove
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
