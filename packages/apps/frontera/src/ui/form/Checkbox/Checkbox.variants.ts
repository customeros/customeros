import { cva } from 'class-variance-authority';

export const CheckboxVariants = cva(
  [
    'flex appearance-none items-center border-solid justify-center border border-gray-300 hover:border-[1px] hover:transition hover:ease-in hover:delay-150 data-[disabled]:cursor-not-allowed data-[disabled]:data-[state=checked]:bg-gray-100 data-[disabled]:data-[state=checked]:border-gray-300 focus:outline-none',
  ],
  {
    variants: {
      size: {
        sm: ['size-4', 'rounded-[4px]'],
        md: ['size-5', 'rounded-[6px]'],
        lg: ['size-6', 'rounded-[8px]'],
        xl: ['size-7', 'rounded-[10px]'],
      },
      colorScheme: {
        primary: [
          'hover:border-primary-600',
          'hover:bg-primary-100',
          'data-[state=checked]:bg-primary-50',
          'data-[state=checked]:border-primary-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-primary-50',
        ],
        gray: [
          'hover:border-gray-600',
          'hover:bg-gray-100',
          'data-[state=checked]:bg-gray-50',
          'data-[state=checked]:border-gray-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-gray-50',
        ],
        warm: [
          'hover:border-warm-600',
          'hover:bg-warm-100',
          'data-[state=checked]:bg-warm-50',
          'data-[state=checked]:border-warm-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-warm-50',
        ],
        error: [
          'hover:border-error-600',
          'hover:bg-error-100',
          'data-[state=checked]:bg-error-50',
          'data-[state=checked]:border-error-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-error-50',
        ],
        rose: [
          'hover:border-rose-600',
          'hover:bg-rose-100',
          'data-[state=checked]:bg-rose-50',
          'data-[state=checked]:border-rose-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-rose-50',
        ],
        warning: [
          'hover:border-warning-600',
          'hover:bg-warning-100',
          'data-[state=checked]:bg-warning-50',
          'data-[state=checked]:border-warning-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-warning-50',
        ],
        blueDark: [
          'hover:border-blueDark-600',
          'hover:bg-blueDark-100',
          'data-[state=checked]:bg-blueDark-50',
          'data-[state=checked]:border-blueDark-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-blueDark-50',
        ],
        teal: [
          'hover:border-teal-600',
          'hover:bg-teal-100',
          'data-[state=checked]:bg-teal-50',
          'data-[state=checked]:border-teal-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-teal-50',
        ],
        success: [
          'hover:border-success-600',
          'hover:bg-success-100',
          'data-[state=checked]:bg-success-50',
          'data-[state=checked]:border-success-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-success-50',
        ],
        moss: [
          'hover:border-moss-600',
          'hover:bg-moss-100',
          'data-[state=checked]:bg-moss-50',
          'data-[state=checked]:border-moss-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-moss-50',
        ],
        greenLight: [
          'hover:border-greenLight-600',
          'hover:bg-greenLight-100',
          'data-[state=checked]:bg-greenLight-50',
          'data-[state=checked]:border-greenLight-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-greenLight-50',
        ],
        violet: [
          'hover:border-violet-600',
          'hover:bg-violet-100',
          'data-[state=checked]:bg-violet-50',
          'data-[state=checked]:border-violet-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-violet-50',
        ],
        fuchsia: [
          'hover:border-fuchsia-600',
          'hover:bg-fuchsia-100',
          'data-[state=checked]:bg-fuchsia-50',
          'data-[state=checked]:border-fuchsia-600',
          'data-[state=checked]:focus:ring-4 data-[state=checked]:focus:ring-fuchsia-50',
        ],
      },
    },
    defaultVariants: {
      size: 'md',
      colorScheme: 'primary',
    },
  },
);
