import { forwardRef, cloneElement, isValidElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

import { cn } from '@ui/utils/cn';
import { useSlots } from '@ui/utils/hooks';

import { Input } from './Input2';

const iconSize = cva([], {
  variants: {
    size: {
      sm: ['size-3 mb-[8px]'],
      md: ['size-4 mb-[5px]'],
      lg: ['size-5'],
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
    <div
      {...props}
      className={twMerge('self-center', className, iconSize({ size }))}
    >
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
      className={twMerge(className, 'self-center', iconSize({ size }))}
    >
      {isValidElement(children) && cloneElement(children, iconProps)}
    </div>
  );
};

interface InputGroupProps extends React.InputHTMLAttributes<HTMLInputElement> {
  border?: boolean;
  children: React.ReactNode;
  ref?: React.Ref<HTMLInputElement>;
}

export const InputGroup = forwardRef<HTMLInputElement, InputGroupProps>(
  ({ border, children }, ref) => {
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
            'flex items-center  w-full border-b gap-1 hover:broder-b hover:border-gray-300 focus-within:hover:border-primary-500 focus-within:border-primary-500 focus-within:border-b hover:transition ease-in-out delay-200',
          )}
        >
          {leftElementSlot}
          {inputSlot}
          {rightElementSlot}
        </div>
      </>
    );
  },
);
