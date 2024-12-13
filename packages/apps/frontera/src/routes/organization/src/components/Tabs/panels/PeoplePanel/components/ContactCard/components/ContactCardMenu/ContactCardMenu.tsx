import { IconButton } from '@ui/form/IconButton';
import { Archive } from '@ui/media/icons/Archive';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';
import { SwitchHorizontal02 } from '@ui/media/icons/SwitchHorizontal02';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

import { ChangeContactOrganizationModal } from './ChangeContactOrganizationModal';

interface ContactCardMenuProps {
  contactId: string;
}

export const ContactCardMenu = ({ contactId }: ContactCardMenuProps) => {
  const { open, onOpen, onClose } = useDisclosure();

  return (
    <>
      <Menu>
        <MenuButton asChild>
          <IconButton
            size='xxs'
            variant='ghost'
            aria-label='More options'
            icon={<DotsVertical className='text-gray-500' />}
          />
        </MenuButton>
        <MenuList>
          <MenuItem className='group/linkedin'>
            <div>
              <LinkedinOutline className='mr-2 text-gray-500 group-hover/linkedin:text-gray-700' />
              <span>Go to LinkedIn profile</span>
            </div>
          </MenuItem>
          <MenuItem className='group/change' onClick={() => onOpen()}>
            <div>
              <SwitchHorizontal02 className='mr-2 text-gray-500 group-hover/change:text-gray-700' />
              <span>Change organization</span>
            </div>
          </MenuItem>
          <MenuItem className='group/archive'>
            <div>
              <Archive className='mr-2 text-gray-500 group-hover/archive:text-gray-700 ' />
              <span>Archive contact</span>
            </div>
          </MenuItem>
        </MenuList>
      </Menu>
      <ChangeContactOrganizationModal
        open={open}
        onClose={onClose}
        contactId={contactId}
      />
    </>
  );
};
