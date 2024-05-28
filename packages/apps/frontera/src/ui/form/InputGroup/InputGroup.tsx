import { twMerge } from 'tailwind-merge';

import { cn } from '@ui/utils/cn';
import { useSlots } from '@ui/utils/hooks';

import { Input } from '../Input';

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
      {children}
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
      {children}
    </div>
  );
};

export interface InputGroupProps extends React.HTMLAttributes<HTMLDivElement> {
  value?: string;
  border?: boolean;
  className?: string;
  children?: React.ReactNode;
}

export const InputGroup = ({
  border,
  children,
  className,
  ...props
}: InputGroupProps) => {
  const [inputSlot, leftElementSlot, rightElementSlot] = useSlots(
    children,
    Input,
    LeftElement,
    RightElement,
  );

  return (
    <>
      <div
        {...props}
        className={cn(
          border ? ' border-gray-200' : 'border-transparent',
          'flex items-center w-full border-b gap-3 hover:broder-b hover:border-gray-300 focus-within:hover:border-primary-500 focus-within:border-primary-500 focus-within:border-b hover:transition ease-in-out delay-200',
          className,
        )}
      >
        {leftElementSlot}
        {inputSlot}
        {rightElementSlot}
      </div>
    </>
  );
};
