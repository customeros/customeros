import { forwardRef } from 'react';
import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

import { cn } from '@ui/utils/cn';

import {
  LeftElement,
  RightElement,
  TextareaGroup,
} from './AutoresizeTextareaGroup';

const sizeTextArea = cva(
  [
    'w-full border-b border-transparent resize-none overflow-hidden gap-3 bg-transparent focus-within:outline-none',
  ],
  {
    variants: {
      size: {
        xs: ['min-h-[19px] text-xs'],
        sm: ['min-h-[19px] text-sm '],
        md: ['min-h-[24px] text-base'],
        lg: ['min-h-[28px] text-lg'],
      },
    },
  },
);

export interface AutoresizeTextareaProps
  extends TextareaAutosizeProps,
    Pick<
      TextareaAutosizeProps,
      'maxRows' | 'minRows' | 'onHeightChange' | 'cacheMeasurements'
    >,
    VariantProps<typeof sizeTextArea> {
  border?: boolean;
  className?: string;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const AutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  AutoresizeTextareaProps
>(({ className, border, leftElement, size, rightElement, ...rest }, ref) => {
  return (
    <TextareaGroup className={cn(className)}>
      {leftElement && <LeftElement>{leftElement}</LeftElement>}
      <ResizeTextarea
        ref={ref}
        {...rest}
        className={twMerge(sizeTextArea({ size }), className)}
      />
      {rightElement && <RightElement>{rightElement}</RightElement>}
    </TextareaGroup>
  );
});
