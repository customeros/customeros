import React from 'react';

import { cn } from '@ui/utils/cn';
import { XCircle } from '@ui/media/icons/XCircle';
import { Action, BilledType } from '@graphql/types';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface ServiceUpdatedActionProps {
  data: Action;
  mode?: 'created' | 'updated' | 'removed';
}

export const ServiceUpdatedAction: React.FC<ServiceUpdatedActionProps> = ({
  data,
  mode = 'updated',
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const metadata = getMetadata(data?.metadata);
  if (!data.content) return null;
  const isTemporary = data.appSource === 'customeros-optimistic-update';
  const formattedContent = formatString(data.content, metadata?.billedType);

  return (
    <div
      className={cn(
        isTemporary
          ? 'opacity-50 cursor-progress'
          : 'opacity-100 cursor-pointer',
        'flex items-center min-h-[40px]',
      )}
      onClick={() => !isTemporary && openModal(data.id)}
    >
      <div className='inline w-[30px]'>
        <FeaturedIcon
          size='md'
          colorScheme={mode === 'created' ? 'primary' : 'gray'}
        >
          {mode === 'removed' ? <XCircle /> : <DotSingle />}
        </FeaturedIcon>
      </div>

      <p className='max-w-[500px] line-clamp-2 ml-2 text-sm text-gray-700 my-1'>
        {formattedContent}
      </p>
    </div>
  );
};

function formatString(str: string, type: string) {
  const digitCount = type === BilledType.Usage ? 4 : 2;
  const regex =
    type === BilledType.Usage ? /\b(\d+\.\d{4})\b/g : /\b(\d+\.\d{2})\b/g;

  return str.replace(regex, (_, number) => {
    return formatCurrency(Number(number), digitCount);
  });
}
