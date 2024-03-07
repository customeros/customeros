import { IconButton } from '@ui/form/IconButton';
import { Portal, useDisclosure } from '@ui/utils';
import { MinusCircle } from '@ui/media/icons/MinusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { CalendarCheck02 } from '@ui/media/icons/CalendarCheck02';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';

interface MilestoneMenuProps {
  opacity?: number;
  transition?: string;
  onRetire?: () => void;
  onAddTask?: () => void;
  onDuplicate?: () => void;
  onSetDueDate?: () => void;
  isMilestoneDone?: boolean;
}

export const MilestoneMenu = ({
  onRetire,
  onAddTask,
  onDuplicate,
  onSetDueDate,
  isMilestoneDone,
  ...buttonProps
}: MilestoneMenuProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Menu
      isOpen={isOpen}
      onClose={onClose}
      onOpen={onOpen}
      placement='bottom-end'
    >
      <MenuButton
        as={IconButton}
        size='xs'
        variant='ghost'
        aria-label='Add milestones'
        icon={<DotsVertical color='gray.400' />}
        {...buttonProps}
        opacity={isOpen ? 1 : buttonProps.opacity}
      />
      <Portal>
        <MenuList minW='10rem'>
          {!isMilestoneDone && (
            <MenuItem
              onClick={onSetDueDate}
              icon={<CalendarCheck02 color='gray.500' />}
            >
              Set due date
            </MenuItem>
          )}
          <MenuItem onClick={onRetire} icon={<MinusCircle color='gray.500' />}>
            Remove
          </MenuItem>
        </MenuList>
      </Portal>
    </Menu>
  );
};
