import { cloneElement, isValidElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

import { cn } from '@ui/utils/cn';
import { useSlots } from '@ui/utils/hooks';

import { Input } from './Input2';

const iconSize = cva([], {
  variants: {
    size: {
      xs: ['w-4', 'h-4'],
      sm: ['w-5', 'h-5'],
      md: ['w-6', 'h-6'],
      lg: ['w-7', 'h-7'],
    },
  },
  defaultVariants: {
    size: 'md',
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
    className: twMerge(iconSize({ size }), className),
    children,
  };

  return (
    <div {...props} className={twMerge(className, iconSize({ size }))}>
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
    className: twMerge(iconSize({ size }), className),
    children,
  };

  return (
    <div
      {...props}
      className={twMerge(
        className,
        'text-gray-500 focus:border-primary-500 focus:border-b',
        iconSize({ size }),
      )}
    >
      {isValidElement(children) && cloneElement(children, iconProps)}
    </div>
  );
};

interface InputGroupProps {
  border?: boolean;
  children: React.ReactNode;
}

export const InputGroup = ({ border, children }: InputGroupProps) => {
  const [inputSlot, leftElementSlot, rightElementSlot] = useSlots(
    children,
    Input,
    LeftElement,
    RightElement,
  );

  return (
    <>
      <div
        className={cn(
          border ? ' border-gray-200' : 'border-transparent',
          'flex items-center w-full border-b gap-3 hover:broder-b hover:border-gray-300 focus-within:hover:border-primary-500 focus-within:border-primary-500 focus-within:border-b hover:transition ease-in-out delay-200',
        )}
      >
        {leftElementSlot}
        {inputSlot}
        {rightElementSlot}
      </div>
    </>
  );
};
