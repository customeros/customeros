import Markdown from 'react-markdown';
import { Children, isValidElement } from 'react';

import { MarkdownEventType } from '@store/TimelineEvents/MarkdownEvent/types';

import { XClose } from '@ui/media/icons/XClose';
import { Link01 } from '@ui/media/icons/Link01';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { getExternalUrl } from '@utils/getExternalLink.ts';
import { CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const GenericActionTypePreview = () => {
  const [_, copy] = useCopyToClipboard();
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as MarkdownEventType;

  return (
    <div className='overflow-hidden rounded-xl pb-6'>
      <CardHeader className='py-4 px-6 pb-1 sticky top-0 rounded-xl bg-white z-[1]'>
        <div className='flex justify-between items-center'>
          <div className='flex items-center'>
            <h2 className='text-base font-semibold '>Website visitor</h2>
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
                  icon={<XClose className='text-gray-500 size-4' />}
                />
              </div>
            </Tooltip>
          </div>
        </div>
      </CardHeader>
      <CardContent className='mt-0 max-h-[calc(100vh-60px-56px)] pt-0 pb-0  text-sm overflow-auto'>
        <Markdown
          className='text-sm whitespace-pre-wrap
        [&>ul]:list-disc [&>ul>li>ul]:list-circle [&>ul>li>ul>li>ul]:list-square
        [&>ol]:list-decimal [&>ol>li>ol]:list-[lower-alpha] [&>ol>li>ol>li>ol]:list-[lower-roman]
        [&>ul>li>ol]:list-decimal [&>ul>li>ol>li>ol]:list-[lower-alpha] [&>ul>li>ol>li>ol>li>ol]:list-[lower-roman]
        [&>ol]:pl-4 [&_ol]:pl-4 [&>ul]:pl-4 [&_ul]:pl-4'
          components={{
            ul: ({ children }) => (
              <ul className={'ml-4 -mt-4 whitespace-nowrap'}>{children}</ul>
            ),
            ol: ({ children }) => (
              <ol className={'ml-4 -mt-4 whitespace-nowrap'}>{children}</ol>
            ),
            li: ({ children }) => {
              const renderContent = () => {
                return Children.map(children, (child, index) => {
                  if (
                    isValidElement(child) &&
                    child.props?.node?.type === 'list'
                  ) {
                    return (
                      <div key={index} className='text-gray-700'>
                        {child}
                      </div>
                    );
                  }

                  return child;
                });
              };

              return <li className='text-gray-700'>{renderContent()}</li>;
            },
            p: ({ children }) => <p className='text-gray-700'>{children}</p>,
            a: ({ children, href }) => (
              <a
                target='_blank'
                rel='noreferrer noopener'
                className='hover:underline'
                href={href ? getExternalUrl(href) : ''}
              >
                {children}
              </a>
            ),
          }}
        >
          {event.content}
        </Markdown>
      </CardContent>
    </div>
  );
};
