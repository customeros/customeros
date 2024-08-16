import { forwardRef } from 'react';
import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import { twMerge } from 'tailwind-merge';

import { InputProps, inputVariants } from '../Input';

interface TextareaProps extends TextareaAutosizeProps {
  size?: InputProps['size'];
  variant?: InputProps['variant'];
  onKeyDown?: (e: React.KeyboardEvent) => void;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ onKeyDown, ...props }, ref) => {
    return (
      <ResizeTextarea
        {...props}
        ref={ref}
        onKeyDown={(e) => {
          if (onKeyDown) onKeyDown(e);
          e.stopPropagation();
        }}
        className={twMerge(
          inputVariants({
            size: props.size,
            variant: props.variant,
          }),
          props.className,
        )}
      />
    );
  },
);
