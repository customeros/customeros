import React, { cloneElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, type VariantProps } from 'class-variance-authority';

import {
  linkButton,
  ghostButton,
  solidButton,
  iconVariant,
  outlineButton,
} from '../Button/Button.variants';

const buttonSize = cva([], {
  variants: {
    size: {
      xs: ['rounded-lg'],
      sm: ['p-1', 'rounded-lg'],
      md: ['p-2.5', 'rounded-lg'],
      lg: ['p-2.5', 'rounded-lg', 'text-base'],
      xl: ['p-3', 'rounded-lg', 'text-base'],
      '2xl': ['p-4', 'gap-3', 'rounded-lg', 'text-lg'],
    },
  },
  defaultVariants: {
    size: 'md',
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
  variant?: 'link' | 'ghost' | 'solid' | 'outline';
}

export const IconButton = ({
  children,
  className,
  colorScheme,
  spinner,
  variant,
  isLoading = false,
  isDisabled = false,
  icon,
  size,
  'aria-label': ariaLabel,
  ...props
}: IconButtonProps) => {
  const buttonVariant = (() => {
    switch (variant) {
      case 'link':
        return linkButton;
      case 'ghost':
        return ghostButton;
      case 'solid':
        return solidButton;
      case 'outline':
        return outlineButton;
      default:
        return solidButton;
    }
  })();

  return (
    <button
      {...props}
      className={twMerge(
        buttonVariant({ colorScheme, className }),
        buttonSize({ className, size }),
        isLoading ? 'opacity-50 cursor-not-allowed' : '',
      )}
      aria-label={ariaLabel}
      disabled={isLoading || isDisabled}
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
};
