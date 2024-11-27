import Markdown from 'react-markdown';

import copy from 'copy-to-clipboard';
import { MarkdownEventType } from '@store/TimelineEvents/MarkdownEvent/types';

import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { CardHeader, CardContent } from '@ui/presentation/Card/Card';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const MarkdownEventPreviewModal = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as MarkdownEventType;

  return (
    <>
      <CardHeader className='py-4 px-6 pb-1 sticky top-0 rounded-xl bg-white z-[1]'>
        <div className='flex justify-between items-center'>
          <div className='flex mb-2 items-center'>
            <h2 className='text-lg font-semibold'>
              Event from {event.markdownEventMetadata.source}
            </h2>
          </div>
          <div className='flex justify-end items-center'>
            <Tooltip side='bottom' label='Copy link to this thread'>
              <div>
                <IconButton
                  size='xs'
                  variant='ghost'
                  color='gray.500'
                  className='mr-1'
                  aria-label='Copy link to this event'
                  onClick={() => copy(window.location.href)}
                  icon={<Link03 className='text-gray-500 size-4' />}
                />
              </div>
            </Tooltip>
            <Tooltip label='Close' side='bottom' aria-label='close'>
              <div>
                <IconButton
                  size='xs'
                  variant='ghost'
                  color='gray.500'
                  onClick={closeModal}
                  aria-label='Close preview'
                  icon={<XClose className='text-gray-500 size-5' />}
                />
              </div>
            </Tooltip>
          </div>
        </div>
      </CardHeader>
      <CardContent className='mt-0 max-h-[calc(100vh-60px-56px)] pb-6 pt-0 text-sm'>
        <Markdown
          className='text-sm'
          components={{
            blockquote: ({ children }) => (
              <blockquote className='text-gray-500 border-l-2 border-gray-300 pl-3 py-1'>
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
              <h1 className='text-base font-bold '>{children}</h1>
            ),
            h2: ({ children }) => (
              <h2 className='text-sm font-medium '>{children}</h2>
            ),
            h3: ({ children }) => (
              <h3 className='text-sm font-medium '>{children}</h3>
            ),
          }}
        >
          {event?.content}
        </Markdown>
      </CardContent>
    </>
  );
};
