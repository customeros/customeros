import { cva } from 'class-variance-authority';

export const CheckboxVariants = cva(
  [
    'flex appearance-none items-center justify-center rounded-[6px] border border-gray-300 hover:border-[1px] hover:transition hover:ease-in hover:delay-150 data-[state=checked]:opacity-100 data-[state=checked]:visible',
  ],
  {
    variants: {
      size: {
        sm: ['size-4'],
        md: ['w-5 h-5'],
        lg: ['size-6'],
        xl: ['size-7'],
      },
      colorScheme: {
        primary: [
          'hover:border-primary-600',
          'hover:bg-primary-100',
          'data-[state=checked]:bg-primary-50',
          'data-[state=checked]:border-primary-600',
        ],
        gray: [],
        warm: [],
        error: [],
        rose: [],
        warning: [],
        blueDark: [],
        teal: [],
        success: [],
        moss: [],
        greenLight: [],
        violet: [],
        fuchsia: [],
      },
    },
    defaultVariants: {
      size: 'md',
      colorScheme: 'primary',
    },
  },
);
