import { forwardRef } from 'react';
import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import { twMerge } from 'tailwind-merge';

import { InputProps, inputVariants } from '../Input';

interface TextareaProps extends TextareaAutosizeProps {
  size?: InputProps['size'];
  variant?: InputProps['variant'];
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  (props, ref) => {
    return (
      <ResizeTextarea
        {...props}
        ref={ref}
        className={twMerge(
          inputVariants({
            size: props?.size,
            variant: props?.variant,
          }),
          props?.className,
        )}
      />
    );
  },
);
