import { cloneElement, isValidElement } from 'react';
import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import { twMerge } from 'tailwind-merge';

import { cn } from '@ui/utils/cn';
import { useSlots } from '@ui/utils/hooks';

interface ElementProps {
  className?: string;
  children: React.ReactNode;
}

export const LeftElement = ({
  children,
  className,
  ...props
}: ElementProps) => {
  return (
    <div {...props} className={twMerge('flex', className)}>
      {isValidElement(children) && cloneElement(children)}
    </div>
  );
};

export const RightElement = ({
  children,
  className,
  ...props
}: ElementProps) => {
  return (
    <div {...props} className={twMerge('flex', className)}>
      {isValidElement(children) && cloneElement(children)}
    </div>
  );
};

export interface TextareaGroupProps
  extends React.HTMLAttributes<HTMLDivElement> {
  border?: boolean;
  className?: string;
  children?: React.ReactNode;
  textareaProps?: TextareaAutosizeProps;
}

/**
 * @deprecated use `<Textarea />` instead
 */
export const TextareaGroup = ({
  border,
  children,
  className,
  textareaProps,
  ...props
}: TextareaGroupProps) => {
  const [textareaSlot, leftElementSlot, rightElementSlot] = useSlots(
    children,
    ResizeTextarea,
    LeftElement,
    RightElement,
  );

  return (
    <>
      <div
        {...props}
        className={cn(
          border ? ' border-gray-200' : 'border-transparent',
          'flex items-center w-full border-b gap-2 py-[2px] mt-1 hover:broder-b hover:border-gray-300 focus-within:hover:border-primary-500 focus-within:border-primary-500 focus-within:border-b hover:transition ease-in-out delay-75',
          className,
        )}
      >
        {leftElementSlot}
        {textareaSlot &&
          cloneElement(textareaSlot as React.ReactElement, {
            ...textareaProps,
          })}
        {rightElementSlot}
      </div>
    </>
  );
};
