import { useDisclosure } from '@ui/utils';
import { File02 } from '@ui/media/icons/File02';
import { File04 } from '@ui/media/icons/File04';
import { IconButton } from '@ui/form/IconButton';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';

interface MasterPlansMenuProps {
  isLoading?: boolean;
  onCreateDefault?: () => void;
  onCreateFromScratch?: () => void;
}

export const MasterPlansMenu = ({
  isLoading,
  onCreateDefault,
  onCreateFromScratch,
}: MasterPlansMenuProps) => {
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
        <MenuItem onClick={onCreateDefault} icon={<File02 color='gray.500' />}>
          Default plan
        </MenuItem>
        <MenuItem
          onClick={onCreateFromScratch}
          icon={<File04 color='gray.500' />}
        >
          From scratch
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
