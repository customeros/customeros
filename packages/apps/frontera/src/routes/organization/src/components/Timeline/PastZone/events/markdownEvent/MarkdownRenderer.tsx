import Markdown from 'react-markdown';
import { Fragment, ReactNode, isValidElement } from 'react';

export const MarkdownRenderer = ({ content }: { content: string }) => {
  return (
    <Markdown
      className='text-sm '
      components={{
        blockquote: ({ children }) => (
          <blockquote className='text-gray-500 border-l border-gray-500 pl-3'>
            {children}
          </blockquote>
        ),
        ul: ({ children }) => (
          <ul className='list-disc list-inside my-1'>{children}</ul>
        ),

        li: ({ children }) => {
          const extractContent = (child: ReactNode) => {
            if (isValidElement(child)) {
              return child.props?.children || child;
            }

            return child;
          };

          // Handle different types of children
          const renderContent = () => {
            if (!children) return null;

            if (
              typeof children === 'string' ||
              typeof children === 'number' ||
              typeof children === 'boolean'
            ) {
              return children;
            }

            if (Array.isArray(children)) {
              return children
                .filter((child) => child !== '\n')
                .map((child, index) => (
                  <Fragment key={index}>{extractContent(child)}</Fragment>
                ));
            }

            return extractContent(children);
          };

          return <li className='list-disc  my-1'>{renderContent()}</li>;
        },
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
      {content}
    </Markdown>
  );
};
