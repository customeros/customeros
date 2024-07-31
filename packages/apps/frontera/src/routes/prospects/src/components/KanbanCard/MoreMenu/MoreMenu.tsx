import { observer } from 'mobx-react-lite';

import { User01 } from '@ui/media/icons/User01';
import { IconButton } from '@ui/form/IconButton';
import { Trash01 } from '@ui/media/icons/Trash01';
import { useStore } from '@shared/hooks/useStore';
import { ArrowsRight } from '@ui/media/icons/ArrowsRight';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

interface MoreMenuProps {
  hasNextSteps: boolean;
  onNextStepsClick: () => void;
}

export const MoreMenu = observer(
  ({ hasNextSteps, onNextStepsClick }: MoreMenuProps) => {
    const store = useStore();

    return (
      <Menu>
        <MenuButton asChild>
          <IconButton
            size='xxs'
            variant='ghost'
            icon={<DotsVertical />}
            aria-label='more options'
          />
        </MenuButton>

        <MenuList>
          <MenuItem onClick={onNextStepsClick}>
            {hasNextSteps ? <Trash01 /> : <ArrowsRight />}
            {hasNextSteps ? 'Remove next step' : 'Add next step'}
          </MenuItem>
          <MenuItem onClick={() => store.ui.commandMenu.toggle('AssignOwner')}>
            <User01 />
            Assign owner
          </MenuItem>
          <MenuItem
            onClick={() => store.ui.commandMenu.toggle('ChangeCurrency')}
          >
            <CurrencyDollarCircle />
            Change currency
          </MenuItem>
        </MenuList>
      </Menu>
    );
  },
);
