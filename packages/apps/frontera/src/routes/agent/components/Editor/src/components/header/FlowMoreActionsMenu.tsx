import { observer } from 'mobx-react-lite';

import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { LayersTwo01 } from '@ui/media/icons/LayersTwo01';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface FlowMoreActionsMenuProps {
  id: string;
}

export const FlowMoreActionsMenu = observer(
  ({ id }: FlowMoreActionsMenuProps) => {
    const store = useStore();

    const handleOpenMenu = (
      type: 'DuplicateFlow' | 'DeleteConfirmationModal',
    ) => {
      store.ui.commandMenu.setContext({
        ...store.ui.commandMenu.context,
        ids: [id],
      });
      store.ui.commandMenu.setType(type);
      store.ui.commandMenu.setOpen(true);
    };

    return (
      <>
        <Menu>
          <MenuButton asChild data-test='flow-editor-more-menu-button'>
            <IconButton
              size='xs'
              aria-label={''}
              variant='ghost'
              icon={<DotsVertical />}
            />
          </MenuButton>
          <MenuList side='bottom' align='center' className='p-0 z-[11]'>
            <MenuItem
              className='flex items-center '
              data-test='duplicate-flow-menu-button'
              onClick={() => {
                handleOpenMenu('DuplicateFlow');
              }}
            >
              <LayersTwo01 className='mr-1 text-gray-500' />
              Duplicate flow
            </MenuItem>
            <MenuItem
              className='flex items-center '
              data-test='duplicate-flow-menu-button'
              onClick={() => {
                handleOpenMenu('DeleteConfirmationModal');
              }}
            >
              <Archive className='mr-1 text-gray-500' />
              Archive flow
            </MenuItem>
          </MenuList>
        </Menu>
      </>
    );
  },
);
