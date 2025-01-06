import { FC } from 'react';

import { MarkdownEventType } from '@store/TimelineEvents/MarkdownEvent/types';

import { cn } from '@ui/utils/cn';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { dataSourceLogo } from '@organization/components/Timeline/PastZone/events/markdownEvent/SourceLogo';
import { MarkdownRenderer } from '@organization/components/Timeline/PastZone/events/markdownEvent/MarkdownRenderer';

import { useTimelineEventPreviewMethodsContext } from '../../../shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const MarkdownEventStub: FC<{ event: MarkdownEventType }> = ({
  event,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();

  return (
    <>
      <Card
        onClick={() => openModal(event.markdownEventMetadata?.id)}
        className={cn(
          'ml-6 shadow-none cursor-pointer text-sm border border-gray-200 bg-white flex max-w-[549px]',
          'rounded-lg hover:shadow-sm transition-all duration-200 ease-out',
        )}
      >
        <CardContent className='p-3 pr-0 overflow-hidden text-sm w-full flex justify-between gap-2'>
          <div>
            <MarkdownRenderer content={event?.content ?? ''} />
          </div>

          <div className='flex items-start min-h-[16px] min-w-[16px] mr-3'>
            {event.markdownEventMetadata?.source &&
              dataSourceLogo?.[event.markdownEventMetadata.source]}
          </div>
        </CardContent>
      </Card>
    </>
  );
};
