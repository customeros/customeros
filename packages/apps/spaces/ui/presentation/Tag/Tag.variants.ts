import { cva } from 'class-variance-authority';

export const tagSubtleVariant = cva(
  ['w-fit', 'flex', 'items-center', 'rounded-md', 'leading-none'],
  {
    variants: {
      colorScheme: {
        primary: ['text-primary-700', 'bg-primary-100'],
        gray: ['text-gray-700', 'bg-gray-100'],
        grayBlue: ['text-grayBlue-700', 'bg-grayBlue-100'],
        warm: ['text-warm-700', 'bg-warm-100'],
        error: ['text-error-700', 'bg-error-100'],
        rose: ['text-rose-700', 'bg-rose-100'],
        warning: ['text-warning-700', 'bg-warning-100'],
        blueDark: ['text-blueDark-700', 'bg-blueDark-100'],
        teal: ['text-teal-700', 'bg-teal-100'],
        success: ['text-success-700', 'bg-success-100'],
        blue: ['text-blue-700', 'bg-blue-100'],
        moss: ['text-moss-700', 'bg-moss-100'],
        greenLight: ['text-greenLight-700', 'bg-greenLight-100'],
        violet: ['text-violet-700', 'bg-violet-100'],
        fuchsia: ['text-fuchsia-700', 'bg-fuchsia-100'],
      },
    },
    defaultVariants: {
      colorScheme: 'gray',
    },
  },
);

export const tagSolidVariant = cva(
  ['w-fit', 'flex', 'items-center', 'rounded-md', 'leading-none'],
  {
    variants: {
      colorScheme: {
        primary: ['bg-primary-500'],
        gray: ['bg-gray-500'],
        grayBlue: ['bg-grayBlue-500'],
        warm: ['bg-warm-500'],
        error: ['bg-error-500'],
        rose: ['bg-rose-500'],
        warning: ['bg-warning-500'],
        blueDark: ['bg-blueDark-500'],
        teal: ['bg-teal-500'],
        success: ['bg-success-500'],
        blue: ['bg-blue-500'],
        moss: ['bg-moss-500'],
        greenLight: ['bg-greenLight-500'],
        violet: ['bg-violet-500'],
        fuchsia: ['bg-fuchsia-500'],
      },
    },
    defaultVariants: {
      colorScheme: 'gray',
    },
  },
);

export const tagOutlineVariant = cva(
  ['w-fit', 'flex', 'items-center', 'rounded-md', 'leading-none'],
  {
    variants: {
      colorScheme: {
        primary: [
          'bg-primary-50',
          'text-primary-700',
          'border',
          'border-solid',
          'border-primary-200',
        ],
        gray: [
          'bg-gray-50',
          'text-gray-700',
          'border',
          'border-solid',
          'border-gray-200',
        ],
        grayBlue: [
          'bg-grayBlue-50',
          'text-grayBlue-700',
          'border',
          'border-solid',
          'border-grayBlue-200',
        ],
        warm: [
          'bg-warm-50',
          'text-warm-700',
          'border',
          'border-solid',
          'border-warm-200',
        ],
        error: [
          'bg-error-50',
          'text-error-700',
          'border',
          'border-solid',
          'border-error-200',
        ],
        rose: [
          'bg-rose-50',
          'text-rose-700',
          'border',
          'border-solid',
          'border-rose-200',
        ],
        warning: [
          'bg-warning-50',
          'text-warning-700',
          'border',
          'border-solid',
          'border-warning-200',
        ],
        blueDark: [
          'bg-blueDark-50',
          'text-blueDark-700',
          'border',
          'border-solid',
          'border-blueDark-200',
        ],
        teal: [
          'bg-teal-50',
          'text-teal-700',
          'border',
          'border-solid',
          'border-teal-200',
        ],
        success: [
          'bg-success-50',
          'text-success-700',
          'border',
          'border-solid',
          'border-success-200',
        ],
        blue: [
          'bg-blue-50',
          'text-blue-700',
          'border',
          'border-solid',
          'border-blue-200',
        ],
        moss: [
          'bg-moss-50',
          'text-moss-700',
          'border',
          'border-solid',
          'border-moss-200',
        ],
        greenLight: [
          'bg-greenLight-50',
          'text-greenLight-700',
          'border',
          'border-solid',
          'border-greenLight-200',
        ],
        violet: [
          'bg-violet-50',
          'text-violet-700',
          'border',
          'border-solid',
          'border-violet-200',
        ],
        fuchsia: [
          'bg-fuchsia-50',
          'text-fuchsia-700',
          'border',
          'border-solid',
          'border-fuchsia-200',
        ],
      },
    },
    defaultVariants: {
      colorScheme: 'gray',
    },
  },
);

export const tagSizeVariant = cva('', {
  variants: {
    size: {
      sm: 'px-2 text-xs',
      md: 'px-2 text-sm',
      lg: 'px-3 text-base',
    },
  },
  defaultVariants: {
    size: 'md',
  },
});
