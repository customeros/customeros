import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

const inputVariants = cva(
  ['w-full', 'ease-in-out', 'delay-200', 'hover:transition'],
  {
    variants: {
      size: {
        xs: ['h-6'],
        sm: ['h-8'],
        md: ['h-10'],
        lg: ['h-12'],
      },
      variant: {
        flushed: [
          'text-gray-700',
          'bg-transparent',
          'placeholder-gray-400',
          'border-b',
          'border-transparent',
          'hover:broder-b',
          'hover:border-gray-300',
          'focus:outline-none',
          'focus:border-b',
          'focus:hover:border-primary-500',
          'focus:border-primary-500',
          'invalid:border-error-500',
        ],
        group: [
          'text-gray-700',
          'bg-transparent',
          'placeholder-gray-400',
          'focus:outline-none',
        ],
      },
    },
    defaultVariants: {
      size: 'md',
      variant: 'flushed',
    },
  },
);

interface InputProps
  extends VariantProps<typeof inputVariants>,
    Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  className?: string;
  placeholder?: string;
}

export const Input = ({ size, variant, className, ...props }: InputProps) => {
  return (
    <input
      {...props}
      className={twMerge(inputVariants({ className, size, variant }))}
    />
  );
};
