import { cva } from 'class-variance-authority';

export const featureIconVariant = cva(
  [
    'flex',
    'justify-center',
    'items-center',
    'rounded-full',
    'overflow-visible',
  ],
  {
    variants: {
      colorScheme: {
        primary: [],
        gray: [],
        grayBlue: [],
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
        blue: [],
        yellow: [],
        orangeDark: [],
      },
    },
    compoundVariants: [
      {
        colorScheme: 'primary',
        className: [
          'bg-primary-100 ring-primary-50 ring-offset-primary-100 text-primary-600',
        ],
      },
      {
        colorScheme: 'gray',
        className: [
          'bg-gray-100 ring-gray-50 ring-offset-gray-100 text-gray-600',
        ],
      },
      {
        colorScheme: 'grayBlue',
        className: [
          'bg-grayBlue-100 ring-grayBlue-50 ring-offset-grayBlue-100 text-grayBlue-600',
        ],
      },
      {
        colorScheme: 'warm',
        className: [
          'bg-warm-100 ring-warm-50 ring-offset-warm-100 text-warm-600',
        ],
      },
      {
        colorScheme: 'error',
        className: [
          'bg-error-100 ring-error-50 ring-offset-error-100 text-error-600',
        ],
      },
      {
        colorScheme: 'rose',
        className: [
          'bg-rose-100 ring-rose-50 ring-offset-rose-100 text-rose-600',
        ],
      },
      {
        colorScheme: 'warning',
        className: [
          'bg-warning-100 ring-warning-50 ring-offset-warning-100 text-warning-600',
        ],
      },
      {
        colorScheme: 'yellow',
        className: [
          'bg-yellow-100 ring-yellow-50 ring-offset-yellow-100 text-yellow-600',
        ],
      },
      {
        colorScheme: 'blueDark',
        className: [
          'bg-blueDark-100 ring-blueDark-50 ring-offset-blueDark-100 text-blueDark-600',
        ],
      },
      {
        colorScheme: 'teal',
        className: [
          'bg-teal-100 ring-teal-50 ring-offset-teal-100 text-teal-600',
        ],
      },
      {
        colorScheme: 'success',
        className: [
          'bg-success-100 ring-success-50 ring-offset-success-100 text-success-600',
        ],
      },
      {
        colorScheme: 'blue',
        className: [
          'bg-blue-100 ring-blue-50 ring-offset-blue-100 text-blue-600',
        ],
      },
      {
        colorScheme: 'moss',
        className: [
          'bg-moss-100 ring-moss-50 ring-offset-moss-100 text-moss-600',
        ],
      },
      {
        colorScheme: 'greenLight',
        className: [
          'bg-greenLight-100 ring-greenLight-50 ring-offset-greenLight-100 text-greenLight-600',
        ],
      },
      {
        colorScheme: 'violet',
        className: [
          'bg-violet-100 ring-violet-50 ring-offset-violet-100 text-violet-600',
        ],
      },
      {
        colorScheme: 'fuchsia',
        className: [
          'bg-fuchsia-100 ring-fuchsia-50 ring-offset-fuchsia-100 text-fuchsia-600',
        ],
      },
      {
        colorScheme: 'orangeDark',
        className: [
          'bg-orangeDark-100 ring-orangeDark-50 ring-offset-orangeDark-100 text-orangeDark-600',
        ],
      },
    ],
  },
);
