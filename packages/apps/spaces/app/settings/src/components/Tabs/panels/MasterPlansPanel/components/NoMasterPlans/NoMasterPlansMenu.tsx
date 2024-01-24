import { Button } from '@ui/form/Button';
import { useDisclosure } from '@ui/utils';
import { File02 } from '@ui/media/icons/File02';
import { File04 } from '@ui/media/icons/File04';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';

interface NoMasterPlansMenuProps {
  isLoading?: boolean;
  onCreateDefault?: () => void;
  onCreateFromScratch?: () => void;
}

export const NoMasterPlansMenu = ({
  isLoading,
  onCreateDefault,
  onCreateFromScratch,
}: NoMasterPlansMenuProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Menu isOpen={isOpen} onClose={onClose} onOpen={onOpen} placement='bottom'>
      <MenuButton
        as={Button}
        variant='outline'
        isLoading={isLoading}
        colorScheme='primary'
      >
        Create plan
      </MenuButton>
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
