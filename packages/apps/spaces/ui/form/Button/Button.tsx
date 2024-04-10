import React, { forwardRef, cloneElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, type VariantProps } from 'class-variance-authority';

import {
  linkButton,
  ghostButton,
  solidButton,
  iconVariant,
  outlineButton,
} from './Button.variants';

export const buttonSize = cva([], {
  variants: {
    size: {
      xs: ['px-2', 'rounded-lg', 'text-xs'],
      sm: ['px-3', 'py-1', 'rounded-lg'],
      md: ['px-4', 'py-[7px]', 'rounded-lg'],
      lg: ['px-[1.125rem]', 'py-2.5', 'rounded-lg', 'text-base'],
      xl: ['px-5', 'py-3', 'rounded-lg', 'text-base'],
      '2xl': ['px-7', 'py-4', 'gap-3', 'rounded-lg', 'text-lg'],
    },
  },
  defaultVariants: {
    size: 'sm',
  },
});

export interface ButtonProps
  extends React.HTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof solidButton>,
    VariantProps<typeof buttonSize> {
  asChild?: boolean;
  isLoading?: boolean;
  loadingText?: string;
  isDisabled?: boolean;
  spinner?: React.ReactElement;
  leftIcon?: React.ReactElement;
  rightIcon?: React.ReactElement;
  variant?: 'link' | 'ghost' | 'solid' | 'outline';
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      leftIcon,
      children,
      className,
      rightIcon,
      colorScheme = 'gray',
      spinner,
      variant = 'outline',
      isLoading = false,
      isDisabled = false,
      size,
      loadingText,
      ...props
    },
    ref,
  ) => {
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
        ref={ref}
        {...props}
        className={twMerge(
          buttonVariant({ colorScheme, className }),
          buttonSize({ className, size }),
          isLoading ? 'opacity-50 cursor-not-allowed' : '',
        )}
        disabled={isLoading || isDisabled}
      >
        {isLoading && spinner && (
          <span className='flex gap-1 relative '>
            {spinner}
            {loadingText}
          </span>
        )}

        {!isLoading && leftIcon && (
          <>
            {cloneElement(leftIcon, {
              className: twMerge(
                iconVariant({
                  size,
                  variant,
                  colorScheme,
                  className: leftIcon.props.className,
                }),
              ),
            })}
          </>
        )}

        {!isLoading && children}
        {!isLoading && rightIcon && (
          <>
            {cloneElement(rightIcon, {
              className: twMerge(
                iconVariant({
                  size,
                  variant,
                  colorScheme,
                  className: rightIcon.props.className,
                }),
              ),
            })}
          </>
        )}
      </button>
    );
  },
);
