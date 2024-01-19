import { useDisclosure } from '@ui/utils';
import { Copy03 } from '@ui/media/icons/Copy03';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/media/icons/Divider';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { SunSetting02 } from '@ui/media/icons/SunSetting02';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';

interface MilestoneMenuProps {
  opacity?: number;
  transition?: string;
  onRetire?: () => void;
  onDuplicate?: () => void;
  onMakeOptional?: () => void;
}

export const MilestoneMenu = ({
  onRetire,
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
        aria-label='Milestone Options'
        icon={<DotsVertical color='gray.400' />}
        {...buttonProps}
        opacity={isOpen ? 1 : buttonProps.opacity}
      />
      <MenuList>
        <MenuItem onClick={onMakeOptional} icon={<Divider color='gray.500' />}>
          Make optional
        </MenuItem>
        <MenuItem onClick={onDuplicate} icon={<Copy03 color='gray.500' />}>
          Duplicate
        </MenuItem>
        <MenuItem onClick={onRetire} icon={<SunSetting02 color='gray.500' />}>
          Retire
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
