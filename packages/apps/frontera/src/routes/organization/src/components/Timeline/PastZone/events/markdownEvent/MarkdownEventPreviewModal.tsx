import copy from 'copy-to-clipboard';
import { MarkdownEventType } from '@store/TimelineEvents/MarkdownEvent/types';

import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { MarkdownRenderer } from '@organization/components/Timeline/PastZone/events/markdownEvent/MarkdownRenderer';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const MarkdownEventPreviewModal = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as MarkdownEventType;

  return (
    <div className='overflow-hidden rounded-xl pb-6'>
      <CardHeader className='py-4 px-6 pb-1 sticky top-0 rounded-xl bg-white z-[1]'>
        <div className='flex justify-between items-center'>
          <div className='flex mb-2 items-center'>
            <h2 className='text-lg font-semibold capitalize'>
              Event from
              <div className='capitalize ml-1 inline-flex'>
                {event.markdownEventMetadata.source.toLowerCase()}
              </div>
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
      <CardContent className='mt-0 max-h-[calc(100vh-60px-56px)] pt-0 pb-0 text-sm overflow-auto'>
        <MarkdownRenderer content={event?.content ?? ''} />
      </CardContent>
    </div>
  );
};
