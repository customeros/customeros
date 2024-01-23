import { useDisclosure } from '@ui/utils';
import { File02 } from '@ui/media/icons/File02';
import { IconButton } from '@ui/form/IconButton';
import { SunSetting02 } from '@ui/media/icons/SunSetting02';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';

interface MasterPlanMenuProps {
  isLoading?: boolean;
  onRetire?: () => void;
  onDuplicate?: () => void;
}

export const MasterPlanMenu = ({
  onRetire,
  isLoading,
  onDuplicate,
}: MasterPlanMenuProps) => {
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
        isLoading={isLoading}
        aria-label='Milestone Options'
        icon={<DotsVertical color='gray.400' />}
      />
      <MenuList minW='10rem'>
        <MenuItem onClick={onDuplicate} icon={<File02 color='gray.500' />}>
          Duplicate
        </MenuItem>
        <MenuItem onClick={onRetire} icon={<SunSetting02 color='gray.500' />}>
          Retire plan
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
