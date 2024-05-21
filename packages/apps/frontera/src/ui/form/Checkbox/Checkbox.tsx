import React, { forwardRef } from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixCheckbox from '@radix-ui/react-checkbox';
import { cva, VariantProps } from 'class-variance-authority';

import { CheckboxVariants } from './Checkbox.variants';

const iconColor = cva(
  ['transition duration-300 ease-in-out transform scale-100'],
  {
    variants: {
      iconColorScheme: {
        primary: ['fill-primary-600'],
        gray: ['fill-gray-600'],
        warm: ['fill-warm-600'],
        error: ['fill-error-600'],
        rose: ['fill-rose-600'],
        warning: ['fill-warning-600'],
        blueDark: ['fill-blueDark-600'],
        teal: ['fill-teal-600'],
        success: ['fill-success-600'],
        moss: ['fill-moss-600'],
        greenLight: ['fill-greenLight-600'],
        violet: ['fill-violet-600'],
        fuchsia: ['fill-fuchsia-600'],
      },
      iconSize: {
        sm: ['size-2'],
        md: ['size-3'],
        lg: ['size-4'],
        xl: ['size-5'],
      },
    },
    defaultVariants: {
      iconColorScheme: 'primary',
      iconSize: 'md',
    },
  },
);

export interface CheckboxProps
  extends VariantProps<typeof CheckboxVariants>,
    VariantProps<typeof iconColor> {
  id?: string;
  disabled?: boolean;
  className?: string;
  defaultChecked?: boolean;
  children?: React.ReactNode;
  isChecked?: boolean | RadixCheckbox.CheckedState;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
  onChange?: (checked: RadixCheckbox.CheckedState | undefined) => void;
}

export const Checkbox = forwardRef<HTMLButtonElement, CheckboxProps>(
  (
    {
      isChecked,
      disabled,
      className,
      size,
      colorScheme,
      iconColorScheme,
      id,
      defaultChecked,
      iconSize,
      onChange,
      children,
      labelProps,
      ...props
    },
    ref,
  ) => {
    return (
      <div className='flex items-center gap-2'>
        <RadixCheckbox.Root
          ref={ref}
          {...props}
          className={twMerge(
            className,
            CheckboxVariants({ size, colorScheme }),
          )}
          defaultChecked={defaultChecked}
          checked={isChecked}
          disabled={disabled}
          onCheckedChange={onChange}
          id={id}
        >
          <RadixCheckbox.Indicator>
            <CheckIcon
              className={twMerge(
                iconColor({
                  iconSize,
                  iconColorScheme,
                  class: disabled ? 'fill-gray-300' : undefined,
                }),
              )}
            />
          </RadixCheckbox.Indicator>
        </RadixCheckbox.Root>
        {children && (
          <label {...labelProps} htmlFor={id}>
            {children}
          </label>
        )}
      </div>
    );
  },
);

const CheckIcon = ({
  className,
  ...props
}: React.SVGAttributes<SVGElement>) => {
  return (
    <svg
      className={className}
      viewBox='0 0 9 9'
      fill='none'
      xmlns='http://www.w3.org/2000/svg'
      {...props}
    >
      <path
        fillRule='evenodd'
        clipRule='evenodd'
        d='M8.53547 0.62293C8.88226 0.849446 8.97976 1.3142 8.75325 1.66099L4.5083 8.1599C4.38833 8.34356 4.19397 8.4655 3.9764 8.49358C3.75883 8.52167 3.53987 8.45309 3.3772 8.30591L0.616113 5.80777C0.308959 5.52987 0.285246 5.05559 0.563148 4.74844C0.84105 4.44128 1.31533 4.41757 1.62249 4.69547L3.73256 6.60459L7.49741 0.840706C7.72393 0.493916 8.18868 0.396414 8.53547 0.62293Z'
      />
    </svg>
  );
};
