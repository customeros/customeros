import React from 'react';

import { cn } from '@ui/utils/cn';
import { FeaturedIcon } from '@ui/media/Icon';
import { File02 } from '@ui/media/icons/File02';
import { Action, ActionType } from '@graphql/types';
import { FileCheck02 } from '@ui/media/icons/FileCheck02';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { FileAttachment02 } from '@ui/media/icons/FileAttachment02';

interface InvoiceStatusChangeActionProps {
  data: Action;
  mode:
    | ActionType.InvoiceIssued
    | ActionType.InvoiceSent
    | ActionType.InvoicePaid
    | ActionType.InvoiceVoided;
}

const iconMap: Record<string, JSX.Element> = {
  [ActionType.InvoiceVoided]: <SlashCircle01 className='text-gray-500' />,
  [ActionType.InvoicePaid]: <FileCheck02 className='text-success-600' />,
  [ActionType.InvoiceSent]: <FileAttachment02 className='text-primary-600' />,
  [ActionType.InvoiceIssued]: <File02 className='text-primary-600' />,
};

const colorSchemeMap: Record<string, string> = {
  [ActionType.InvoiceVoided]: 'gray',
  [ActionType.InvoicePaid]: 'success',
  [ActionType.InvoiceSent]: 'primary',
  [ActionType.InvoiceIssued]: 'primary',
};

const InvoiceStatusChangeAction: React.FC<InvoiceStatusChangeActionProps> = ({
  data,
  mode,
}) => {
  const isTemporary = data.appSource === 'customeros-optimistic-update';

  if (!data.content) return <div>No data available</div>;

  const formattedContent = formatInvoiceText(data.content);

  return (
    <div
      className={cn('flex items-center pointer ', {
        'not-allowed': isTemporary,
        'opacity-50': isTemporary,
      })}
    >
      <FeaturedIcon size='md' minW='10' colorScheme={colorSchemeMap[mode]}>
        {iconMap[mode]}
      </FeaturedIcon>
      <p className='my-1 max-w-[500px] ml-2 text-sm text-gray-700 '>
        {formattedContent}
      </p>
    </div>
  );
};

const formatInvoiceText = (text: string) => {
  const invoiceNumberPattern = /N° \w+-\d+/;
  const amountPattern = /\$\d{1,3}(,\d{3})*(\.\d{2})?/;

  const invoiceNumberMatch = text.match(invoiceNumberPattern);
  const amountMatch = text.match(amountPattern);

  const invoiceNumber = invoiceNumberMatch ? invoiceNumberMatch[0] : '';
  const amount = amountMatch ? amountMatch[0] : '';

  const beforeInvoiceNumber = text.split(invoiceNumberPattern)[0];
  const betweenInvoiceNumberAndAmount = text
    .split(invoiceNumberPattern)[1]
    .split(amountPattern)[0];
  const afterAmount = text.split(amountPattern)[1];

  return (
    <div>
      {beforeInvoiceNumber}
      <span className='font-medium'>{invoiceNumber}</span>
      {betweenInvoiceNumberAndAmount}
      <span className='font-medium'>{amount}</span>
      {afterAmount}
    </div>
  );
};

export default InvoiceStatusChangeAction;
