import React from 'react';

import { cn } from '@ui/utils/cn';
import { Action } from '@graphql/types';
import { FeaturedIcon } from '@ui/media/Icon';
import { File02 } from '@ui/media/icons/File02';
import { FileCheck02 } from '@ui/media/icons/FileCheck02';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { FileAttachment02 } from '@ui/media/icons/FileAttachment02';
import { getMetadata } from '@organization/src/components/Timeline/PastZone/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface InvoiceStatusChangeActionProps {
  data: Action;
}

export enum ExtendedInvoiceStatus {
  Void = 'Void',
  Paid = 'Paid',
  // Adding new statuses
  Send = 'SEND',
  Issued = 'ISSUED',
  Unknown = 'UNKNOWN',
}

const iconMap: Record<string, JSX.Element> = {
  [ExtendedInvoiceStatus.Void]: <SlashCircle01 className='text-gray-500' />,
  [ExtendedInvoiceStatus.Paid]: <FileCheck02 className='text-success-600' />,
  [ExtendedInvoiceStatus.Send]: (
    <FileAttachment02 className='text-primary-600' />
  ),
  [ExtendedInvoiceStatus.Issued]: <File02 className='text-primary-600' />,
};

const colorSchemeMap: Record<string, string> = {
  [ExtendedInvoiceStatus.Void]: 'gray',
  [ExtendedInvoiceStatus.Paid]: 'success',
  [ExtendedInvoiceStatus.Send]: 'primary',
  [ExtendedInvoiceStatus.Issued]: 'primary',
};

const InvoiceStatusChangeAction: React.FC<InvoiceStatusChangeActionProps> = ({
  data,
}) => {
  const { handleOpenInvoice } = useTimelineEventPreviewMethodsContext();
  const metadata = getMetadata(data?.metadata);
  const isTemporary = data.appSource === 'customeros-optimistic-update';

  if (!data.content) return <div>No data available</div>;

  const formattedContent = formatInvoiceText(data.content);
  const status: ExtendedInvoiceStatus = data.content.includes('issued')
    ? ExtendedInvoiceStatus.Issued
    : data.content.includes('Sent')
    ? ExtendedInvoiceStatus.Send
    : (metadata?.status as ExtendedInvoiceStatus) ??
      ExtendedInvoiceStatus.Unknown;

  return (
    <div
      role='button'
      tabIndex={0}
      onClick={() => !isTemporary && handleOpenInvoice(data.id)}
      className={cn('flex items-center pointer ', {
        'not-allowed': isTemporary,
        'opacity-50': isTemporary,
      })}
    >
      <FeaturedIcon size='md' minW='10' colorScheme={colorSchemeMap[status]}>
        {iconMap[status]}
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
