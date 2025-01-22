import Markdown from 'react-markdown';
import { Children, isValidElement } from 'react';

import { Action } from '@graphql/types';
import { Globe06 } from '@ui/media/icons/Globe06';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface GenericActionTypeProps {
  data: Action;
}

export const GenericActionType = ({ data }: GenericActionTypeProps) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();

  if (!data?.content) return null;

  return (
    <div>
      <div
        onClick={() => openModal(data.id)}
        className='flex cursor-pointer min-h-[40px]'
      >
        <Globe06 className='text-gray-500 mt-0.5 mr-2' />
        <Markdown
          className='text-sm whitespace-pre-line
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

            a: ({ children }) => (
              <a href={''} className={'pointer-events-none'}>
                {children}
              </a>
            ),
          }}
        >
          {data.content}
        </Markdown>
      </div>
    </div>
  );
};
