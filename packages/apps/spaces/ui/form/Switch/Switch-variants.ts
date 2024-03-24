import { cva } from 'class-variance-authority';

export const switchVariants = cva(
  ['bg-gray-100', 'outline-none', 'cursor-default', 'relative', 'rounded-full'],
  {
    variants: {
      size: {
        sm: ['w-[25px] h-[16px]'],
        md: ['w-[34px] h-[20px]'],
        lg: ['w-[50px] h-[25px]'],
      },
      colorScheme: {
        primary: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-primary-50',
          'data-[state=checked]:bg-primary-600',
        ],
        gray: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-gray-50',
          'data-[state=checked]:bg-gray-500',
        ],
        warm: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-warm-50',
          'data-[state=checked]:bg-warm-500',
        ],
        error: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-error-50',
          'data-[state=checked]:bg-error-500',
        ],
        rose: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-rose-50',
          'data-[state=checked]:bg-rose-500',
        ],
        warning: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-warning-50',
          'data-[state=checked]:bg-warning-500',
        ],
        blueDark: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-blueDark-50',
          'data-[state=checked]:bg-blueDark-500',
        ],
        teal: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-teal-50',
          'data-[state=checked]:bg-teal-500',
        ],
        success: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-success-50',
          'data-[state=checked]:bg-success-500',
        ],
        moss: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-moss-50',
          'data-[state=checked]:bg-moss-500',
        ],
        greenLight: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-greenLight-50',
          'data-[state=checked]:bg-greenLight-500',
        ],
        violet: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-violet-50',
          'data-[state=checked]:bg-violet-500',
        ],
        fuchsia: [
          'hover:bg-gray-200',
          'focus:ring-4',
          'focus:ring-fuchsia-50',
          'data-[state=checked]:bg-fuchsia-500',
        ],
      },
    },
    defaultVariants: {
      size: 'md',
      colorScheme: 'primary',
    },
  },
);
