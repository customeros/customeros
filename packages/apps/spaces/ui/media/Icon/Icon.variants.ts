import { cva } from 'class-variance-authority';

export const featureIconVariant = cva(
  ['flex', 'justify-center', 'items-center , rounded-full'],
  {
    variants: {
      colorScheme: {
        primary: [],
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
    compoundVariants: [
      {
        colorScheme: 'primary',
        className: [
          'bg-primary-100 ring-primary-50 ring-offset-primary-100 text-primary-500',
        ],
      },
      {
        colorScheme: 'gray',
        className: [
          'bg-gray-100 ring-gray-50 ring-offset-gray-100 text-gray-500',
        ],
      },
      {
        colorScheme: 'grayBlue',
        className: [
          'bg-grayBlue-100 ring-grayBlue-50 ring-offset-grayBlue-100 text-grayBlue-500',
        ],
      },
      {
        colorScheme: 'warm',
        className: [
          'bg-warm-100 ring-warm-50 ring-offset-warm-100 text-warm-500',
        ],
      },
      {
        colorScheme: 'error',
        className: [
          'bg-error-100 ring-error-50 ring-offset-error-100 text-error-500',
        ],
      },
      {
        colorScheme: 'rose',
        className: [
          'bg-rose-100 ring-rose-50 ring-offset-rose-100 text-rose-500',
        ],
      },
      {
        colorScheme: 'warning',
        className: [
          'bg-warning-100 ring-warning-50 ring-offset-warning-100 text-warning-500',
        ],
      },
      {
        colorScheme: 'blueDark',
        className: [
          'bg-blueDark-100 ring-blueDark-50 ring-offset-blueDark-100 text-blueDark-500',
        ],
      },
      {
        colorScheme: 'teal',
        className: [
          'bg-teal-100 ring-teal-50 ring-offset-teal-100 text-teal-500',
        ],
      },
      {
        colorScheme: 'success',
        className: [
          'bg-success-100 ring-success-50 ring-offset-success-100 text-success-500',
        ],
      },
      {
        colorScheme: 'moss',
        className: [
          'bg-moss-100 ring-moss-50 ring-offset-moss-100 text-moss-500',
        ],
      },
      {
        colorScheme: 'greenLight',
        className: [
          'bg-greenLight-100 ring-greenLight-50 ring-offset-greenLight-100 text-greenLight-500',
        ],
      },
      {
        colorScheme: 'violet',
        className: [
          'bg-violet-100 ring-violet-50 ring-offset-violet-100 text-violet-500',
        ],
      },
      {
        colorScheme: 'fuchsia',
        className: [
          'bg-fuchsia-100 ring-fuchsia-50 ring-offset-fuchsia-100 text-fuchsia-500',
        ],
      },
    ],
  },
);
