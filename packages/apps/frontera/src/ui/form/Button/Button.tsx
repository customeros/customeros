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
      xxs: ['py-0.5', 'px-1.5', 'rounded-md', 'text-xs'],
      xs: ['px-2', 'py-1', 'rounded-md', 'text-sm'],
      sm: ['px-3', 'py-2', 'rounded-lg', 'text-sm'],
      md: ['px-4', 'py-[10px]', 'rounded-lg', 'text-sm'],
      lg: ['px-[1.125rem]', 'py-2.5', 'rounded-lg', 'text-base'],
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
  dataTest?: string;
  isLoading?: boolean;
  isDisabled?: boolean;
  loadingText?: string;
  leftIcon?: React.ReactElement;
  rightIcon?: React.ReactElement;
  leftSpinner?: React.ReactElement;
  rightSpinner?: React.ReactElement;
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
      leftSpinner,
      rightSpinner,
      loadingText,
      variant = 'outline',
      isLoading = false,
      isDisabled = false,
      size = 'sm',
      dataTest,
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
        data-test={dataTest}
        style={{ outline: 'none' }}
        disabled={isLoading || isDisabled}
        className={twMerge(
          buttonVariant({ colorScheme, className }),
          buttonSize({ className, size }),
          isLoading ? 'opacity-50 cursor-not-allowed' : '',
        )}
      >
        {leftIcon && (
          <>
            {isLoading && leftSpinner
              ? leftSpinner
              : cloneElement(leftIcon, {
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

        {isLoading ? loadingText || children : children}
        {rightIcon && (
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

        {isLoading && rightSpinner && (
          <span className='flex gap-1 relative items-center'>
            {rightSpinner}
          </span>
        )}
      </button>
    );
  },
);
