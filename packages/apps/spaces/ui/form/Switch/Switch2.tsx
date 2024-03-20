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

interface SwitchProps
  extends Omit<RadixSwitch.SwitchProps, 'onChange'>,
    VariantProps<typeof switchVariants> {
  className?: string;
  isChecked?: boolean;
  isDisabled?: boolean;
  isRequired?: boolean;
  onChange?: (value: boolean) => void;
}

export const Switch = ({
  colorScheme,
  isDisabled,
  isRequired,
  isChecked,
  className,
  onChange,
  size,
  ...props
}: SwitchProps) => {
  return (
    <RadixSwitch.Root
      onCheckedChange={onChange}
      checked={isChecked}
      required={isRequired}
      disabled={isDisabled}
      className={twMerge(switchVariants({ colorScheme, size }), className)}
      style={
        {
          WebkitTapHighlightColor: 'rgba(0, 0, 0, 0)',
        } as React.CSSProperties
      }
    >
      <RadixSwitch.Thumb className={twMerge(thumbSizes({ size }), className)} />
    </RadixSwitch.Root>
  );
};
