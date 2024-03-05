import { twMerge } from 'tailwind-merge';
import { cva, type VariantProps } from 'class-variance-authority';

const button = cva(
  [
    'inline-flex',
    'items-center',
    'justify-center',
    'whitespace-nowrap',
    'gap-2',
    'text-sm',
    'font-semibold',
    'shadow-xs',
    'outline-none',
    'transition',
    'disabled:pointer-events-none',
    'disabled:opacity-50',
  ],
  {
    variants: {
      isDestructive: {
        true: [],
      },
      variant: {
        primary: [
          'text-white',
          'border',
          'border-solid',
          'bg-primary-600',
          'hover:bg-primary-700',
          'focus:bg-primary-700',
          'border-primary-600',
          'hover:border-primary-700',
          'focus:shadow-ringPrimary',
          'focus-visible:shadow-ringPrimary',
        ],
        secondary: [
          'text-gray-600',
          'bg-white',
          'border',
          'border-solid',
          'border-gray-300',
          'hover:bg-gray-50',
          'hover:text-gray-700',
          'focus:bg-gray-50',
          'focus:shadow-ringSecondary',
          'focus-visible:shadow-ringSecondary',
        ],
        secondaryAccent: [
          'text-primary-600',
          'bg-white',
          'border',
          'border-solid',
          'border-primary-300',
          'hover:bg-primary-50',
          'hover:text-primary-700',
          'focus:bg-primary-50',
          'focus:shadow-ringPrimary',
          'focus-visible:shadow-ringPrimary',
        ],
        tertiary: [
          'bg-transparent',
          'text-gray-500',
          'shadow-none',
          'hover:text-gray-700',
          'focus:text-gray-700',
          'hover:bg-gray-50',
          'focus:bg-gray-50',
        ],
        tertiaryAccent: [
          'bg-transparent',
          'text-primary-500',
          'shadow-none',
          'hover:text-primary-700',
          'focus:text-primary-700',
          'hover:bg-primary-50',
          'focus:bg-primary-50',
        ],
      },
      size: {
        sm: ['px-3.5', 'py-2', 'rounded-lg'],
        md: ['px-4', 'py-2.5', 'rounded-lg'],
        lg: ['px-[1.125rem]', 'py-2.5', 'rounded-lg', 'text-base'],
        xl: ['px-5', 'py-3', 'rounded-lg', 'text-base'],
        '2xl': ['px-7', 'py-4', 'gap-3', 'rounded-lg', 'text-lg'],
      },
    },
    compoundVariants: [
      {
        variant: 'primary',
        isDestructive: true,
        className: [
          'text-white',
          'border',
          'border-solid',
          'bg-error-600',
          'hover:bg-error-700',
          'focus:bg-error-700',
          'border-error-600',
          'hover:border-error-700',
          'focus:shadow-ringDestructive',
          'focus-visible:shadow-ringDestructive',
        ],
      },
      {
        variant: 'secondary',
        isDestructive: true,
        className: [
          'text-error-600',
          'bg-white',
          'border',
          'border-solid',
          'border-error-300',
          'hover:bg-error-50',
          'hover:text-error-700',
          'focus:bg-error-50',
          'focus:shadow-ringDestructive',
          'focus-visible:shadow-ringDestructive',
        ],
      },
      {
        variant: 'secondaryAccent',
        isDestructive: true,
        className: [
          'text-error-600',
          'bg-white',
          'border',
          'border-solid',
          'border-error-300',
          'hover:bg-error-50',
          'hover:text-error-700',
          'focus:bg-error-50',
          'focus:shadow-ringDestructive',
          'focus-visible:shadow-ringDestructive',
        ],
      },
      {
        variant: 'tertiary',
        isDestructive: true,
        className: [
          'bg-transparent',
          'text-error-500',
          'shadow-none',
          'hover:text-error-700',
          'focus:text-error-700',
          'hover:bg-error-50',
          'focus:bg-error-50',
        ],
      },
      {
        variant: 'tertiaryAccent',
        isDestructive: true,
        className: [
          'bg-transparent',
          'text-error-500',
          'shadow-none',
          'hover:text-error-700',
          'focus:text-error-700',
          'hover:bg-error-50',
          'focus:bg-error-50',
        ],
      },
    ],
    defaultVariants: {
      variant: 'primary',
      size: 'md',
    },
  },
);

interface ButtonProps
  extends React.HTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof button> {
  asChild?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Button = ({
  size,
  variant,
  leftIcon,
  children,
  className,
  rightIcon,
  isDestructive,
  ...props
}: ButtonProps) => {
  return (
    <button
      {...props}
      className={twMerge(button({ size, variant, isDestructive, className }))}
    >
      {leftIcon && <span>{leftIcon}</span>}
      {children}
      {rightIcon && <span>{leftIcon}</span>}
    </button>
  );
};
