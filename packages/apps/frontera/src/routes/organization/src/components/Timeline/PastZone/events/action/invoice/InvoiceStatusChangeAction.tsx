import React, { useMemo } from 'react';

import { cn } from '@ui/utils/cn';
import { File02 } from '@ui/media/icons/File02';
import { Action, ActionType } from '@graphql/types';
import { FileX02 } from '@ui/media/icons/FileX02.tsx';
import { FileCheck02 } from '@ui/media/icons/FileCheck02';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { FileAttachment02 } from '@ui/media/icons/FileAttachment02';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface InvoiceStatusChangeActionProps {
  data: Action;
  mode:
    | ActionType.InvoiceIssued
    | ActionType.InvoiceSent
    | ActionType.InvoicePaid
    | ActionType.InvoiceOverdue
    | ActionType.InvoiceVoided;
}

interface InvoiceStubMetadata {
  id: string;
  status: string;
  number: string;
  amount: number;
  currency: string;
}

const iconMap: Record<string, JSX.Element> = {
  [ActionType.InvoiceVoided]: <SlashCircle01 className='text-gray-500' />,
  [ActionType.InvoicePaid]: <FileCheck02 className='text-success-600' />,
  [ActionType.InvoiceSent]: <FileAttachment02 className='text-primary-600' />,
  [ActionType.InvoiceIssued]: <File02 className='text-primary-600' />,
  [ActionType.InvoiceOverdue]: <FileX02 className='text-primary-600' />,
};

const colorSchemeMap: Record<
  string,
  | 'primary'
  | 'gray'
  | 'grayBlue'
  | 'warm'
  | 'error'
  | 'rose'
  | 'warning'
  | 'blueDark'
  | 'teal'
  | 'success'
  | 'blue'
  | 'moss'
  | 'greenLight'
  | 'violet'
  | 'fuchsia'
> = {
  [ActionType.InvoiceVoided]: 'gray',
  [ActionType.InvoicePaid]: 'success',
  [ActionType.InvoiceSent]: 'primary',
  [ActionType.InvoiceOverdue]: 'warning',
  [ActionType.InvoiceIssued]: 'primary',
};

const InvoiceStatusChangeAction: React.FC<InvoiceStatusChangeActionProps> = ({
  data,
  mode,
}) => {
  const isTemporary = data.appSource === 'customeros-optimistic-update';
  const { handleOpenInvoice } = useTimelineEventPreviewMethodsContext();

  const metadata = useMemo(() => {
    return getMetadata(data?.metadata) as unknown as InvoiceStubMetadata;
  }, [data?.metadata]);
  if (!data.content) return <div>No data available</div>;

  const formattedContent = formatInvoiceText(data.content, metadata);

  return (
    <div
      role='button'
      tabIndex={0}
      onClick={() =>
        !isTemporary && metadata?.id && handleOpenInvoice(metadata.id)
      }
      className={cn('flex items-center pointer focus:outline-none min-h-10 ', {
        'not-allowed': isTemporary || !metadata?.id,
        'opacity-50': isTemporary,
      })}
    >
      <FeaturedIcon
        size='md'
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        colorScheme={colorSchemeMap[mode] as any}
        className='mr-[10px]'
      >
        {iconMap[mode]}
      </FeaturedIcon>
      <p className='my-1 max-w-[500px] ml-5 text-sm text-gray-700 '>
        {formattedContent}
      </p>
    </div>
  );
};

const formatInvoiceText = (text: string, metadata: InvoiceStubMetadata) => {
  if (!metadata) {
    return text;
  }
  const invoiceNumberPattern = /N° \w+-\d+/;
  const formattedAmount = formatCurrency(metadata.amount, 2, metadata.currency);
  const beforeInvoiceNumber = text.split(invoiceNumberPattern)[0];
  const betweenInvoiceNumberAndAmount = text
    .split(invoiceNumberPattern)[1]
    .replace(formattedAmount, '')
    .trim();

  const afterAmount = text.split(`${metadata.amount}`)[1];

  return (
    <div>
      {beforeInvoiceNumber}
      <span className='font-medium mr-1'>N° {metadata.number}</span>
      {betweenInvoiceNumberAndAmount}
      <span className='font-medium mx-1'>{formattedAmount}</span>
      {afterAmount}
    </div>
  );
};

export default InvoiceStatusChangeAction;
