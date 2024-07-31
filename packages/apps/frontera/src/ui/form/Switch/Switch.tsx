import { forwardRef } from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixSwitch from '@radix-ui/react-switch';
import { cva, VariantProps } from 'class-variance-authority';

import { switchVariants } from './Switch-variants';

const thumbSizes = cva(
  [
    'bg-white',
    'rounded-full',
    'block',
    'transition-transform duration-100 translate-x-0.5 will-change-transform',
  ],
  {
    variants: {
      size: {
        sm: [],
        md: [],
        lg: [],
      },
    },
    compoundVariants: [
      {
        size: 'sm',
        className: 'size-3 data-[state=checked]:translate-x-[11px]',
      },
      {
        size: 'md',
        className: 'size-4 data-[state=checked]:translate-x-[16px]',
      },
      {
        size: 'lg',
        className: 'size-5 data-[state=checked]:translate-x-[28px]',
      },
    ],
    defaultVariants: {
      size: 'md',
    },
  },
);

export interface SwitchProps
  extends Omit<RadixSwitch.SwitchProps, 'onChange'>,
    VariantProps<typeof switchVariants> {
  className?: string;
  isChecked?: boolean;
  isInvalid?: boolean;

  isDisabled?: boolean;
  isRequired?: boolean;
  onChange?: (value: boolean) => void;
}

export const Switch = forwardRef<HTMLButtonElement, SwitchProps>(
  (
    {
      size,
      colorScheme,
      className,
      isChecked,
      isInvalid,
      isDisabled,
      isRequired,
      onChange,
      ...rest
    },
    ref,
  ) => {
    const invalidContainer =
      isInvalid && ' data-[state=checked]:bg-warning-500';
    const invalidThumb =
      isInvalid &&
      'after:content-["!"] after:absolute after:top-[-2px] after:left-0 after:right-0 after:text-xs after:text-warning-500 font-bold';

    return (
      <RadixSwitch.Root
        ref={ref}
        checked={isChecked}
        required={isRequired}
        disabled={isDisabled}
        onCheckedChange={onChange}
        style={
          {
            WebkitTapHighlightColor: 'rgba(0, 0, 0, 0)',
          } as React.CSSProperties
        }
        className={twMerge(
          switchVariants({ colorScheme, size }),
          className,
          invalidContainer,
        )}
        {...rest}
      >
        <RadixSwitch.Thumb
          className={twMerge(thumbSizes({ size }), className, invalidThumb)}
        />
      </RadixSwitch.Root>
    );
  },
);
