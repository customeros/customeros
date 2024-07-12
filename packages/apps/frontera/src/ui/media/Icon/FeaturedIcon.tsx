import { ReactNode, cloneElement, isValidElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';

import { featureIconVariant } from './Icon.variants';

const sizeIcon = cva([], {
  variants: {
    size: {
      sm: ['size-4 ring-[4px] ring-offset-[6px]'],
      md: ['size-5 ring-[6px] ring-offset-[7px]'],
      lg: ['size-6 ring-[8px] ring-offset-[8px]'],
      xl: ['size-7 ring-[10px] ring-offset-[9px]'],
    },
  },
  defaultVariants: {
    size: 'md',
  },
});

export interface FeaturedIconStyleProps
  extends VariantProps<typeof featureIconVariant>,
    VariantProps<typeof sizeIcon> {
  className?: string;
  children:
    | ReactNode
    | ((props: React.HTMLAttributes<HTMLDivElement>) => ReactNode);
}

export const FeaturedIcon = ({
  children,
  colorScheme = 'primary',
  size = 'md',
  className,
}: FeaturedIconStyleProps) => {
  const iconProps: React.HTMLAttributes<HTMLDivElement> = {
    ...featureIconVariant,
    className: twMerge(
      sizeIcon({ size }),
      featureIconVariant({ className, colorScheme }),
      className,
    ),
  };

  return isValidElement(children)
    ? cloneElement(children, iconProps)
    : typeof children === 'function'
    ? // eslint-disable-next-line @typescript-eslint/no-explicit-any
      children(iconProps as any)
    : children;
};
