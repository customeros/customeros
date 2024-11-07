import { cloneElement } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Clock } from '@ui/media/icons/Clock';
import { InvoiceStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { renderStatusNode } from '@shared/components/Invoice/Cells';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

export const PaymentStatusSelect = observer(
  ({
    invoiceNumber,
    variant = 'invoice-finder',
  }: {
    invoiceNumber: string;
    variant?: 'invoice-preview' | 'invoice-finder';
  }) => {
    const store = useStore();

    const invoice = invoiceNumber
      ? store.invoices.value.get(invoiceNumber)
      : null;

    const invoiceStatus = invoice?.value?.status;
    const Status = renderStatusNode(invoiceStatus) ?? <>{invoiceStatus}</>;
    const isPaid = invoiceStatus === InvoiceStatus.Paid;

    if (!invoice) return;

    return (
      <Menu>
        <MenuButton
          asChild
          onClick={(e) => e.stopPropagation()}
          disabled={invoiceStatus === InvoiceStatus.Scheduled}
        >
          {cloneElement(Status, {
            className: cn(
              'cursor-pointer',
              variant === 'invoice-preview' && 'px-2 py-0.5 text-md',
              invoiceStatus === InvoiceStatus.Scheduled &&
                'opacity-50 cursor-not-allowed',
            ),
          })}
        </MenuButton>
        <MenuList
          side='bottom'
          align='center'
          onCloseAutoFocus={(e) => e.preventDefault()}
        >
          {invoiceStatus !== InvoiceStatus.Void && (
            <MenuItem
              disabled={isPaid}
              onClick={() => {
                invoice.value.status = InvoiceStatus.Void;
                invoice?.commit();
              }}
            >
              <div className='flex gap-2 items-center'>
                <SlashCircle01
                  className={cn(isPaid ? 'text-gray-400' : 'text-gray-500')}
                />
                <span>Void</span>
              </div>
            </MenuItem>
          )}
          {invoiceStatus !== InvoiceStatus.Paid && (
            <MenuItem
              onClick={() => {
                invoice.value.status = InvoiceStatus.Paid;
                invoice?.commit();
              }}
            >
              <div className='flex gap-2 items-center'>
                <CheckCircle className='text-gray-500' />
                <span>Paid</span>
              </div>
            </MenuItem>
          )}
          {invoiceStatus !== InvoiceStatus.Due && (
            <MenuItem
              onClick={() => {
                invoice.value.status = InvoiceStatus.Due;
                invoice?.commit();
              }}
            >
              <div className='flex gap-2 items-center'>
                <Clock className='text-gray-500' />
                <span>Due</span>
              </div>
            </MenuItem>
          )}
        </MenuList>
      </Menu>
    );
  },
);
