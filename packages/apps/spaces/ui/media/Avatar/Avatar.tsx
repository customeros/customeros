import React from 'react';
import { cloneElement } from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixAvatar from '@radix-ui/react-avatar';
import { AvatarImageProps } from '@radix-ui/react-avatar';
import { cva, VariantProps } from 'class-variance-authority';

const avatarBadgeSize = cva([], {
  variants: {
    badgeSize: {
      xs: ['w-[6px] h-[6px]'],
      sm: ['w-[8px] h-[8px]'],
      md: ['w-[10px] h-[10px]'],
      lg: ['w-[12px] h-[12px]'],
      xl: ['w-[14px] h-[14px]'],
      '2xl': ['w-[16px] h-[16px]'],
    },
    borderRadius: {
      xs: ['ring-[2px] ring-white'],
      sm: ['ring-[2px] ring-white'],
      md: ['ring-[2px] ring-white'],
      lg: ['ring-[2px] ring-white'],
      xl: ['ring-[2px] ring-white'],
      '2xl': ['border-[8px] ring-white'],
    },
    badgePosition: {
      xs: ['transform -translate-x-[-12px] -translate-y-[-10px]'],
      sm: ['transform -translate-x-[-12px] -translate-y-[-13px]'],
      md: ['transform -translate-x-[-20px] -translate-y-[-15px]'],
      lg: ['transform -translate-x-[-20px] -translate-y-[-20px]'],
      xl: ['transform -translate-x-[-20px] -translate-y-[-25px]'],
      '2xl': ['transform -translate-x-[-25px] -translate-y-[-25px]'],
    },
  },
  compoundVariants: [
    {
      badgeSize: 'xs',
      borderRadius: 'xs',
      badgePosition: 'xs',
    },
    {
      badgeSize: 'sm',
      borderRadius: 'sm',
      badgePosition: 'sm',
    },
    {
      badgeSize: 'md',
      borderRadius: 'md',
      badgePosition: 'md',
    },
    {
      badgeSize: 'lg',
      borderRadius: 'lg',
      badgePosition: 'lg',
    },
    {
      badgeSize: 'xl',
      borderRadius: 'xl',
      badgePosition: 'xl',
    },
    {
      badgeSize: '2xl',
      borderRadius: '2xl',
      badgePosition: '2xl',
    },
  ],
});

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

const textSize = cva([], {
  variants: {
    textSizes: {
      xs: ['text-xs'],
      sm: ['text-sm'],
      md: ['text-base'],
      lg: ['text-lg'],
      xl: ['text-xl'],
      '2xl': ['text-2xl'],
    },
  },
  defaultVariants: {
    textSizes: 'md',
  },
});

interface AvatarProps
  extends VariantProps<typeof avatarStyle>,
    VariantProps<typeof avatarBadgeSize>,
    VariantProps<typeof textSize>,
    AvatarImageProps {
  src?: string;
  name?: string;
  className?: string;
  icon?: React.ReactNode;
  badge?: React.ReactElement;
}

export const Avatar: React.FC<AvatarProps> = ({
  icon,
  name,
  src,
  size,
  textSizes = 'md',
  variant,
  badgeSize,
  className,
  color,
  badge,
  ...props
}: AvatarProps) => {
  const emptyFallbackWords = name?.split(' ');

  let emptyFallbackLetters = '';
  if (!emptyFallbackWords) return null;
  if (emptyFallbackWords.length > 1) {
    emptyFallbackLetters = `${emptyFallbackWords[0]?.[0]}${
      emptyFallbackWords[emptyFallbackWords.length - 1]?.[0]
    }`.toLocaleUpperCase();
  } else {
    emptyFallbackLetters = emptyFallbackWords[0]?.[0].toLocaleUpperCase();
  }

  return (
    <RadixAvatar.Root
      className={twMerge(avatarStyle({ size, variant, className }))}
    >
      {src && (
        <RadixAvatar.Image
          {...props}
          className={'h-full w-full relative rounded-[inherit] object-cover'}
          src={src}
        />
      )}
      {icon && !name && !src && (
        <RadixAvatar.Fallback
          className={twMerge(
            'leading-1 flex h-full w-full items-center justify-center font-medium',
            textSize({ textSizes }),
          )}
        >
          {icon}
        </RadixAvatar.Fallback>
      )}
      {(!icon || name) && !src && (
        <RadixAvatar.Fallback
          className={twMerge(
            'leading-1 flex h-full w-full items-center justify-center font-bold',
            textSize({ textSizes }),
          )}
        >
          {emptyFallbackLetters}
        </RadixAvatar.Fallback>
      )}
      {badge &&
        cloneElement(badge, {
          className: twMerge(
            avatarBadgeSize({
              badgeSize: size,
              badgePosition: size,
              borderRadius: size,
            }),
            badge.props.className,
          ),
        })}
    </RadixAvatar.Root>
  );
};

interface AvatarBadgeProps extends VariantProps<typeof avatarBadgeSize> {
  className?: string;
}

export const AvatarBadge: React.FC<AvatarBadgeProps> = ({
  className,
  badgePosition,
  badgeSize,
}: AvatarBadgeProps) => {
  return (
    <div
      className={twMerge([
        className,
        'rounded-full absolute ',
        avatarBadgeSize({ badgeSize, badgePosition: badgeSize }),
      ])}
    />
  );
};
