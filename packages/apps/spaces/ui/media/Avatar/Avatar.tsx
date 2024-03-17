import { twMerge } from 'tailwind-merge';
import * as RadixAvatar from '@radix-ui/react-avatar';
import { AvatarImageProps } from '@radix-ui/react-avatar';
import { cva, VariantProps } from 'class-variance-authority';

const avatarStyle = cva(
  [
    'bg-blackA1',
    'inline-flex',
    'select-none',
    'items-center',
    'justify-center',
    'overflow-hidden',
    'align-middle',
  ],
  {
    variants: {
      variant: {
        circle: ['rounded-full', 'bg-primary-100', 'text-primary-700'],
        shadowed: [
          'bg-primary-100',
          'text-primary-700',
          'rounded-full',
          'ring-offset-1',
          'ring-4',
          'ring-offset-primary-200/10',
          'ring-primary-100/50',
        ],
        roundedSquareSmall: ['rounded-sm', 'text-primary-700', 'bg-primary-50'],
        roundedSquare: ['rounded-md', 'text-primary-700', 'bg-primary-50'],
        roundedSquareShadowed: [
          'rounded-lg',
          'bg-primary-100',
          'text-primary-700',
          'ring-offset-1',
          'ring-4',
          'ring-offset-primary-200/10',
          'ring-primary-100/50',
        ],
        outline: [
          'rounded-full',
          'bg-primary-50',
          'border',
          'border-primary-200',
        ],
      },

      size: {
        xs: ['w-6 h-6'],
        sm: ['w-8 h-8'],
        md: ['w-10 h-10'],
        lg: ['w-12 h-12'],
        xl: ['w-14 h-14'],
        '2xl': ['w-16 h-16'],
      },
    },

    defaultVariants: {
      size: 'lg',
      variant: 'circle',
    },
  },
);
interface AvatarDemoProps
  extends VariantProps<typeof avatarStyle>,
    AvatarImageProps {
  src?: string;
  name?: string;
  className?: string;
  icon?: React.ReactNode;
}

export const AvatarDemo = ({
  icon,
  name,
  src,
  size,
  variant,
  className,
  color,
  ...props
}: AvatarDemoProps) => {
  const emptyFallbackLetters = name?.split(' ').map((word) => word[0]);

  return (
    <RadixAvatar.Root
      className={twMerge(avatarStyle({ size, variant, className }))}
    >
      {src && (
        <RadixAvatar.Image
          {...props}
          className={'h-full w-full rounded-[inherit] object-cover'}
          src={src}
        />
      )}
      {icon && !name && !src && (
        <RadixAvatar.Fallback
          delayMs={600}
          className=' leading-1 flex h-full w-full items-center justify-center text-[15px] font-medium'
        >
          {icon}
        </RadixAvatar.Fallback>
      )}
      {!icon && !src && name && (
        <RadixAvatar.Fallback
          className='leading-1 flex h-full w-full items-center justify-center text-[15px] font-semibold'
          delayMs={600}
        >
          {emptyFallbackLetters}
        </RadixAvatar.Fallback>
      )}
    </RadixAvatar.Root>
  );
};
