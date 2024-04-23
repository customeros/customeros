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
    'w-full border-b border-transparent placeholder-gray-400 leading-6 resize-none overflow-hidden gap-2 bg-transparent focus-within:outline-none',
  ],
  {
    variants: {
      size: {
        xs: ['min-h-6 h-6'],
        sm: ['min-h-8 h-8'],
        md: ['min-h-10 h-10'],
        lg: ['min-h-12 h-12'],
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
  label?: string;
  border?: boolean;
  className?: string;
  size?: 'xs' | 'sm' | 'md';
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
  labelProps?: React.HTMLAttributes<HTMLLabelElement>;
}

export const AutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  AutoresizeTextareaProps
>(
  (
    {
      className,
      border,
      leftElement,
      label,
      labelProps,
      size,
      rightElement,
      ...rest
    },
    ref,
  ) => {
    return (
      <div>
        {label && <label {...labelProps}>{label}</label>}
        <TextareaGroup className={cn(className)}>
          {leftElement && <LeftElement>{leftElement}</LeftElement>}
          <ResizeTextarea
            ref={ref}
            {...rest}
            className={twMerge(sizeTextArea({ size }), className)}
          />
          {rightElement && <RightElement>{rightElement}</RightElement>}
        </TextareaGroup>
      </div>
    );
  },
);
