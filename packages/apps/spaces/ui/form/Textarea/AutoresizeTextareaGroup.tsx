import { cloneElement, isValidElement } from 'react';
import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

import { cn } from '@ui/utils/cn';
import { useSlots } from '@ui/utils/hooks';

const iconSize = cva([], {
  variants: {
    size: {
      xs: ['size-4'],
      sm: ['size-5'],
      md: ['size-6'],
    },
  },
  defaultVariants: {
    size: 'sm',
  },
});

interface ElementProps extends VariantProps<typeof iconSize> {
  className?: string;
  children: React.ReactNode;
}

export const LeftElement = ({
  children,
  size,
  className,
  ...props
}: ElementProps) => {
  const iconProps = {
    ...iconSize,
    ...props,
    className: twMerge(className, iconSize({ size })),
    children,
  };

  return (
    <div {...props} className={twMerge('flex', className, iconSize({ size }))}>
      {isValidElement(children) && cloneElement(children, iconProps)}
    </div>
  );
};

export const RightElement = ({
  children,
  size,
  className,
  ...props
}: ElementProps) => {
  const iconProps = {
    ...iconSize,
    ...props,
    className: twMerge(className, iconSize({ size })),
    children,
  };

  return (
    <div {...props} className={twMerge('flex', className, iconSize({ size }))}>
      {isValidElement(children) && cloneElement(children, iconProps)}
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
