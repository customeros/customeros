import { useRef, useState, useEffect } from 'react';

import { cn } from '@ui/utils/cn.ts';

interface TruncatedTextProps extends React.HTMLAttributes<HTMLDivElement> {
  text: string;
  maxLines?: number;
}

export const TruncatedText = ({
  text,
  maxLines = 5,
  className,
  ...rest
}: TruncatedTextProps) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [showButton, setShowButton] = useState(false);
  const textRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (textRef.current) {
      const lineHeight = parseInt(
        getComputedStyle(textRef.current).lineHeight,
        10,
      );
      const maxHeight = lineHeight * maxLines;

      setShowButton(textRef.current.scrollHeight > maxHeight);
    }
  }, [text, maxLines]);

  const toggleExpand = () => {
    setIsExpanded(!isExpanded);
  };

  return (
    <div className='w-full'>
      <div
        ref={textRef}
        className={cn('overflow-hidden', className ? className : '')}
        style={{
          display: '-webkit-box',
          lineClamp: isExpanded ? 'none' : maxLines,
          WebkitLineClamp: isExpanded ? 'none' : maxLines,
          WebkitBoxOrient: 'vertical',
        }}
        {...rest}
      >
        {text}
      </div>
      {showButton && (
        <button
          onClick={toggleExpand}
          className={cn(
            'inline-block text-primary-700 underline',
            className ? className : '',
          )}
        >
          {isExpanded ? 'Show less' : 'Show more'}
        </button>
      )}
    </div>
  );
};
