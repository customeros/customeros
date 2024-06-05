import { Clock } from '@ui/media/icons/Clock';
import { InvoiceStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { renderStatusNode } from '@shared/components/Invoice/Cells';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const StatusMenuButton = ({
  status,
  id,
}: {
  id: string;
  status?: InvoiceStatus | null;
}) => {
  const store = useStore();

  const invoice = store.invoices?.value?.get(id);

  const handleUpdateStatus = (newStatus: InvoiceStatus) => {
    invoice?.update((prev) => ({
      ...prev,
      status: newStatus,
    }));
  };

  return (
    <Menu>
      <MenuButton aria-label='Status'>{renderStatusNode(status)}</MenuButton>
      <MenuList align='start' side='bottom' className='w-[200px] shadow-xl'>
        {status !== InvoiceStatus.Paid && (
          <MenuItem
            color='gray.700'
            onClick={() => handleUpdateStatus(InvoiceStatus.Paid)}
          >
            <CheckCircle className='text-gray-500 mr-2' />
            Paid
          </MenuItem>
        )}
        {status !== InvoiceStatus.Void && (
          <MenuItem
            color='gray.700'
            onClick={() => handleUpdateStatus(InvoiceStatus.Void)}
          >
            <SlashCircle01 className='text-gray-500 mr-2' />
            Void
          </MenuItem>
        )}

        {status !== InvoiceStatus.Due && (
          <MenuItem
            color='gray.700'
            onClick={() => handleUpdateStatus(InvoiceStatus.Due)}
          >
            <Clock className='text-gray-500 mr-2' />
            Due
          </MenuItem>
        )}
      </MenuList>
    </Menu>
  );
};
