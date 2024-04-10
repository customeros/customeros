import { forwardRef } from 'react';
import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import { cn } from '@ui/utils/cn';

import {
  LeftElement,
  RightElement,
  TextareaGroup,
} from './AutoresizeTextareaGroup';

export interface AutoresizeTextareaProps
  extends TextareaAutosizeProps,
    Pick<
      TextareaAutosizeProps,
      'maxRows' | 'minRows' | 'onHeightChange' | 'cacheMeasurements'
    > {
  border?: boolean;
  className?: string;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const AutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  AutoresizeTextareaProps
>(({ className, border, leftElement, rightElement, ...rest }, ref) => {
  return (
    <TextareaGroup className={cn(className)}>
      {leftElement && <LeftElement>{leftElement}</LeftElement>}
      <ResizeTextarea
        ref={ref}
        {...rest}
        className='h-[24px] w-full border-b border-transparent resize-none overflow-hidden gap-3 bg-transparent focus-within:outline-none'
      />
      {rightElement && <RightElement>{rightElement}</RightElement>}
    </TextareaGroup>
  );
});
