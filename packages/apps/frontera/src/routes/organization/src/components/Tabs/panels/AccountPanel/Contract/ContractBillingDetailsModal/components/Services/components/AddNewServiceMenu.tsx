import { observer } from 'mobx-react-lite';

import { BilledType } from '@graphql/types';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { RefreshCcw02 } from '@ui/media/icons/RefreshCcw02.tsx';
import { CalendarDate } from '@ui/media/icons/CalendarDate.tsx';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

interface AddNewServiceMenuProps {
  isInline?: boolean;
  contractId: string;
}
//TODO K
export const AddNewServiceMenu = observer(
  ({ isInline, contractId }: AddNewServiceMenuProps) => {
    const store = useStore();
    const contractLineItemsStore = store.contractLineItems;

    return (
      <>
        <Menu>
          <MenuButton asChild>
            {isInline ? (
              <p>Add a service</p>
            ) : (
              <IconButton
                size='xxs'
                icon={<Plus />}
                className='ml-1'
                variant='outline'
                colorScheme='gray'
                aria-label='Add a service'
                data-test='contract-card-add-sli'
              />
            )}
          </MenuButton>

          <MenuList align='end' side='bottom' className='p-0'>
            <MenuItem
              className='flex items-center text-base'
              data-test={'add-new-service-menu-subscription'}
              onClick={() =>
                contractLineItemsStore.create({
                  billingCycle: BilledType.Monthly,
                  contractId,
                  description: 'Unnamed',
                })
              }
            >
              <RefreshCcw02 className='mr-2 text-gray-500' />
              Subscription
            </MenuItem>
            <MenuItem
              className='flex items-center text-base'
              data-test={'add-new-service-menu-one-time'}
              onClick={() =>
                contractLineItemsStore.create({
                  billingCycle: BilledType.Once,
                  contractId,
                  description: 'Unnamed',
                })
              }
            >
              <CalendarDate className='mr-2 text-gray-500' />
              One-time
            </MenuItem>
          </MenuList>
        </Menu>
      </>
    );
  },
);
