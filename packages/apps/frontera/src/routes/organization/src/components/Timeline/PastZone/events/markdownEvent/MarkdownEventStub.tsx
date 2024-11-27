import { FC } from 'react';
import Markdown from 'react-markdown';

import { MarkdownEventType } from '@store/TimelineEvents/MarkdownEvent/types';

import { cn } from '@ui/utils/cn';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { dataSourceLogo } from '@organization/components/Timeline/PastZone/events/markdownEvent/SourceLogo.tsx';

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
          'ml-6 shadow-xs cursor-pointer text-sm border border-gray-200 bg-white flex max-w-[549px]',
          'rounded-lg hover:shadow-md transition-all duration-200 ease-out',
        )}
      >
        <CardContent className='p-3 pr-0 overflow-hidden text-sm flex gap-2 '>
          <div>
            <Markdown
              components={{
                blockquote: ({ children }) => (
                  <blockquote className='text-gray-500 border-l border-gray-500 pl-3'>
                    {children}
                  </blockquote>
                ),
                ul: ({ children }) => (
                  <ul className='list-disc list-inside'>{children}</ul>
                ),
                ol: ({ children }) => (
                  <ul className='list-decimal list-inside'>{children}</ul>
                ),
                h1: ({ children }) => (
                  <h1 className='text-sm font-bold'>{children}</h1>
                ),

                h2: ({ children }) => (
                  <h2 className='text-sm font-medium'>{children}</h2>
                ),
              }}
            >
              {event?.content}
            </Markdown>
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
