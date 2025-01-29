import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';
import { XSquare } from '@ui/media/icons/XSquare';
import { PlayCircle } from '@ui/media/icons/PlayCircle';
import { PauseCircle } from '@ui/media/icons/PauseCircle';
import { BracketsPlus } from '@ui/media/icons/BracketsPlus';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface ProductItemMenuProps {
  id: string;
  closed?: boolean;
  paused?: boolean;
  contractId: string;
  allowPausing?: boolean;
  allowAddModification?: boolean;
  handlePauseService: (paused: boolean) => void;
  handleCloseService: (isClosed: boolean) => void;
}

export const ProductItemMenu = observer(
  ({
    id,
    contractId,
    allowAddModification,
    handleCloseService,
    handlePauseService,
    allowPausing,
    paused,
  }: ProductItemMenuProps) => {
    const store = useStore();
    const contractLineItemsStore = store.contractLineItems;

    return (
      <>
        <Menu>
          <MenuButton
            className={cn(
              `flex items-center max-h-5 p-1 py-2 hover:bg-gray-100 rounded translate-x-2 outline:0`,
            )}
          >
            <DotsVertical className='text-gray-400' />
          </MenuButton>
          <MenuList align='end' side='bottom' className='p-0'>
            {allowAddModification ? (
              <MenuItem
                className='flex items-center text-base'
                onClick={() =>
                  contractLineItemsStore?.create({
                    id,
                    contractId,
                  })
                }
              >
                <BracketsPlus className='mr-2 text-gray-500' />
                Add modification
              </MenuItem>
            ) : (
              <div />
            )}

            {allowPausing ? (
              <MenuItem
                className='flex items-center text-base'
                onClick={() => handlePauseService(!paused)}
              >
                {paused ? (
                  <>
                    <PlayCircle className='mr-2 text-gray-500' />
                    Resume this product
                  </>
                ) : (
                  <>
                    <PauseCircle className='mr-2 text-gray-500' />
                    Pause this product
                  </>
                )}
              </MenuItem>
            ) : (
              <div />
            )}

            <MenuItem
              className='flex items-center text-base'
              onClick={() => handleCloseService(true)}
            >
              <XSquare className='mr-2 text-gray-500' />
              End the service
            </MenuItem>
          </MenuList>
        </Menu>
      </>
    );
  },
);
