import React, { forwardRef, cloneElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, type VariantProps } from 'class-variance-authority';

import { iconVariant } from './IconButton.variants';
import {
  ghostButton,
  solidButton,
  outlineButton,
} from '../Button/Button.variants';

const buttonSize = cva([], {
  variants: {
    size: {
      xxs: ['p-1', 'rounded-md'],
      xs: ['p-1.5', 'rounded-md'],
      sm: ['p-2', 'rounded-lg', 'text-lg'],
      md: ['p-[10px]', 'rounded-lg'],
      lg: ['p-[10px]', 'rounded-lg'],
    },
  },
  defaultVariants: {
    size: 'sm',
  },
});

export interface IconButtonProps
  extends React.HTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof solidButton>,
    VariantProps<typeof buttonSize> {
  asChild?: boolean;
  isLoading?: boolean;
  isDisabled?: boolean;
  'aria-label': string;
  icon: React.ReactElement;
  spinner?: React.ReactElement;
  variant?: 'ghost' | 'solid' | 'outline';
}

export const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(
  (
    {
      children,
      className,
      colorScheme = 'gray',
      spinner,
      variant = 'outline',
      isLoading = false,
      isDisabled = false,
      icon,
      size = 'sm',
      'aria-label': ariaLabel,
      ...props
    },
    ref,
  ) => {
    const buttonVariant = (() => {
      switch (variant) {
        case 'ghost':
          return ghostButton;
        case 'solid':
          return solidButton;
        case 'outline':
          return outlineButton;
        default:
          return outlineButton;
      }
    })();

    return (
      <button
        ref={ref}
        {...props}
        aria-label={ariaLabel}
        disabled={isLoading || isDisabled}
        className={twMerge(
          buttonVariant({ colorScheme, className }),
          buttonSize({ className, size }),
          isLoading ? 'opacity-50 cursor-not-allowed' : '',
        )}
      >
        {isLoading && spinner && (
          <span className='relative inline-flex'>{spinner}</span>
        )}

        {!isLoading && icon && (
          <>
            {cloneElement(icon, {
              className: twMerge(
                iconVariant({
                  size,
                  variant,
                  colorScheme,
                  className: icon.props.className,
                }),
              ),
            })}
          </>
        )}
      </button>
    );
  },
);
