import { forwardRef, KeyboardEvent, InputHTMLAttributes } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

export const inputVariants = cva(
  [
    'w-full',
    'ease-in-out',
    'delay-50',
    'hover:transition',
    'disabled:cursor-not-allowed',
    '[&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none',
  ],
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
          'resize-none',
        ],
        outline: [
          'text-gray-700',
          'bg-white',
          'placeholder-gray-400',
          'border',
          'px-2',
          'rounded-[4px]',
          'border-gray-300',
          'focus:outline-none',
          'focus:border-primary-500',
          'invalid:border-error-500',
        ],
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
    Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
  dataTest?: string;
  className?: string;
  placeholder?: string;
  allowKeyDownEventPropagation?: boolean;
  onKeyDown?: (e: KeyboardEvent<HTMLInputElement>) => void;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      size,
      variant,
      allowKeyDownEventPropagation,
      className,
      onKeyDown,
      dataTest,
      ...rest
    },
    ref,
  ) => {
    return (
      <input
        {...rest}
        ref={ref}
        data-1p-ignore
        data-test={dataTest}
        className={twMerge(inputVariants({ className, size, variant }))}
        onKeyDown={(e) => {
          if (onKeyDown) {
            onKeyDown(e);

            return;
          }

          e.stopPropagation();
        }}
      />
    );
  },
);
