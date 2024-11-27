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
        <CardContent className='p-3 pr-0 overflow-hidden text-sm w-full flex justify-between gap-2'>
          <div>
            <Markdown
              skipHtml
              components={{
                blockquote: ({ children }) => (
                  <blockquote className='text-gray-500 border-l border-gray-500 pl-3'>
                    {children}
                  </blockquote>
                ),
                ul: ({ children }) => (
                  <ul className='list-disc list-inside my-1'>{children}</ul>
                ),
                ol: ({ children }) => (
                  <ol className='list-decimal list-inside my-1'>{children}</ol>
                ),
                h1: ({ children }) => (
                  <h1 className='text-sm font-bold mt-1'>{children}</h1>
                ),

                h2: ({ children }) => (
                  <h2 className='text-sm font-medium mt-1'>{children}</h2>
                ),
                p: ({ children }) => <p className='text-sm my-1'>{children}</p>,
                a: ({ children, href }) => {
                  return (
                    <a href={href} target='_blank' rel='noreferrer noopener'>
                      {children}
                    </a>
                  );
                },
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
