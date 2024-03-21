import { cloneElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

import { cn } from '@ui/utils/cn';

const inputVariants = cva(['w-full'], {
  variants: {
    size: {
      xs: ['h-6'],
      sm: ['h-8'],
      md: ['h-10'],
      lg: ['h-12'],
    },
    variant: {
      flushed: [
        'text-gray-700',
        'bg-transparent',
        'placeholder-gray-400',
        'focus:outline-none',
        'invalid:border-error-500',
      ],
    },
  },
  defaultVariants: {
    size: 'md',
    variant: 'flushed',
  },
});

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

interface InputProps
  extends VariantProps<typeof inputVariants>,
    VariantProps<typeof iconSize>,
    Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  border?: boolean;
  className?: string;
  placeholder?: string;
  leftIcon?: React.ReactElement;
  rightIcon?: React.ReactElement;
  onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

export const Input = ({
  size,
  variant,
  onChange,
  leftIcon,
  className,
  rightIcon,
  placeholder,
  border,
  ...props
}: InputProps) => {
  return (
    <>
      <div
        className={cn(
          border ? ' border-gray-200' : 'border-transparent',
          'flex items-center w-full border-b gap-3 hover:broder-b hover:border-gray-300 focus-within:hover:border-primary-500 focus-within:border-primary-500 focus-within:border-b hover:transition ease-in-out delay-200  ',
        )}
      >
        {leftIcon && (
          <>
            {cloneElement(leftIcon, {
              className: twMerge(
                'text-gray-500  ',
                'focus:border-primary-500 focus:border-b',
                iconSize({ size, className: leftIcon.props.className }),
              ),
            })}
          </>
        )}
        <input
          {...props}
          className={twMerge(inputVariants({ className, size, variant }))}
          onChange={onChange}
          placeholder={placeholder}
        />
        {rightIcon && (
          <>
            {cloneElement(rightIcon, {
              className: twMerge(
                'text-gray-500',
                iconSize({ size, className: rightIcon.props.className }),
              ),
            })}
          </>
        )}
      </div>
    </>
  );
};

//
