import { MarkdownEventType } from '@store/TimelineEvents/MarkdownEvent/types';

import { XClose } from '@ui/media/icons/XClose';
import { Link01 } from '@ui/media/icons/Link01';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { MarkdownRenderer } from '@organization/components/Timeline/PastZone/events/markdownEvent/MarkdownRenderer';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const MarkdownEventPreviewModal = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as MarkdownEventType;
  const [_, copy] = useCopyToClipboard();

  return (
    <div className='overflow-hidden rounded-xl pb-6'>
      <CardHeader className='py-4 px-6 pb-1 sticky top-0 rounded-xl bg-white z-[1]'>
        <div className='flex justify-between items-center'>
          <div className='flex mb-2 items-center'>
            <h2 className='text-base font-semibold capitalize'>
              Event from
              <div className='capitalize ml-1 inline-flex'>
                {event.markdownEventMetadata.source.toLowerCase()}
              </div>
            </h2>
          </div>
          <div className='flex justify-end items-baseline'>
            <Tooltip side='bottom' label='Copy link to this thread'>
              <div>
                <IconButton
                  size='xs'
                  variant='ghost'
                  color='gray.500'
                  className='mr-1'
                  aria-label='Copy link to this event'
                  icon={<Link01 className='text-gray-500 size-4' />}
                  onClick={() => copy(window.location.href, 'Link copied')}
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
