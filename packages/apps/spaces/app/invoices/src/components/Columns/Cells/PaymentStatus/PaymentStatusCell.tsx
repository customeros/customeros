import { cloneElement } from 'react';

import { cn } from '@ui/utils/cn';
import { Clock } from '@ui/media/icons/Clock';
import { InvoiceStatus } from '@graphql/types';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { renderStatusNode } from '@shared/components/Invoice/Cells';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

import { useTableActionState } from '../../../../state/TableActionState.atom';

export const PaymentStatusCell = ({
  value,
  invoiceId,
}: {
  invoiceId: string;
  value: InvoiceStatus | null;
}) => {
  const [_, setTableActionState] = useTableActionState();
  const Status = renderStatusNode(value) ?? <>{value}</>;
  const isPaid = value === InvoiceStatus.Paid;

  const handleClick = (status: InvoiceStatus) => {
    setTableActionState((prev) => ({
      ...prev,
      targetId: invoiceId,
      targetStatus: status,
      isConfirming: status === InvoiceStatus.Void,
    }));
  };

  return (
    <Menu>
      <MenuButton asChild>
        {cloneElement(Status, { className: 'cursor-pointer' })}
      </MenuButton>
      <MenuList
        align='center'
        side='bottom'
        onCloseAutoFocus={(e) => e.preventDefault()}
      >
        <MenuItem
          disabled={isPaid}
          onClick={() => handleClick(InvoiceStatus.Void)}
        >
          <div className='flex gap-2 items-center'>
            <SlashCircle01
              className={cn(isPaid ? 'text-gray-400' : 'text-gray-500')}
            />
            <span>Void</span>
          </div>
        </MenuItem>
        <MenuItem onClick={() => handleClick(InvoiceStatus.Paid)}>
          <div className='flex gap-2 items-center'>
            <CheckCircle className='text-gray-500' />
            <span>Paid</span>
          </div>
        </MenuItem>
        <MenuItem onClick={() => handleClick(InvoiceStatus.Due)}>
          <div className='flex gap-2 items-center'>
            <Clock className='text-gray-500' />
            <span>Due</span>
          </div>
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
