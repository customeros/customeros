import React, { forwardRef } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

export const inputVariants = cva(
  ['w-full', 'ease-in-out', 'delay-50', 'hover:transition'],
  {
    variants: {
      size: {
        xs: ['min-h-6 text-sm'],
        sm: ['min-h-8 text-sm'],
        md: ['min-h-10 text-base'],
        lg: ['min-h-12 text-lg'],
      },
      variant: {
        flushed: [
          'text-gray-700',
          'bg-transparent',
          'placeholder-gray-400',
          'border-b',
          'border-transparent',
          'hover:broder-b',
          'hover:border-gray-300',
          'focus:outline-none',
          'focus:border-b',
          'focus:hover:border-primary-500',
          'focus:border-primary-500',
          'invalid:border-error-500',
        ],
        group: [
          'text-gray-700',
          'bg-transparent',
          'placeholder-gray-400',
          'focus:outline-none',
        ],
        unstyled: [
          'text-gray-700',
          'bg-transparent',
          'placeholder-gray-400',
          'focus:outline-none',
        ],
        outline: [],
      },
    },
    defaultVariants: {
      size: 'md',
      variant: 'flushed',
    },
  },
);

export interface InputProps
  extends VariantProps<typeof inputVariants>,
    Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  className?: string;
  placeholder?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ size, variant, className, ...rest }, ref) => {
    return (
      <input
        {...rest}
        ref={ref}
        className={twMerge(inputVariants({ className, size, variant }))}
      />
    );
  },
);
