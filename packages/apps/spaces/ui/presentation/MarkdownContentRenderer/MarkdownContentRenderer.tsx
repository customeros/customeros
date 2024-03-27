import { FC } from 'react';
import ReactMarkdown from 'react-markdown';

import { twMerge } from 'tailwind-merge';

interface MarkdownContentRendererProps
  extends React.HTMLAttributes<HTMLDivElement> {
  markdownContent: string;
  showAsInlineText?: boolean;
}

export const MarkdownContentRenderer: FC<MarkdownContentRendererProps> = ({
  markdownContent,
  showAsInlineText,
  className,
  ...rest
}) => {
  const textRendererClass = showAsInlineText
    ? 'inline-text-renderer'
    : 'block-text-renderer';

  return (
    <ReactMarkdown
      {...rest}
      className={twMerge(
        'flex flex-col block-text-renderer',
        textRendererClass,
        className,
      )}
    >
      {markdownContent}
    </ReactMarkdown>
  );
};
