import { useMemo, forwardRef, cloneElement, isValidElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { VariantProps } from 'class-variance-authority';

import { cn } from '@ui/utils/cn';
import { useSlots } from '@ui/utils/hooks';
import { XClose } from '@ui/media/icons/XClose';

import {
  tagSizeVariant,
  tagSolidVariant,
  tagSubtleVariant,
  tagOutlineVariant,
} from './Tag.variants';

const allVariants = {
  solid: tagSolidVariant,
  subtle: tagSubtleVariant,
  outline: tagOutlineVariant,
};

export interface TagProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof tagSizeVariant>,
    VariantProps<typeof tagSolidVariant> {
  variant?: 'subtle' | 'solid' | 'outline';
}

export const Tag = forwardRef<HTMLDivElement, TagProps>(
  (
    { size, children, className, colorScheme, variant = 'outline', ...props },
    ref,
  ) => {
    const [leftIconSlot, rightIconSlot, labelSlot, closeButtonSlot] = useSlots(
      children,
      TagLeftIcon,
      TagRightIcon,
      TagLabel,
      TagCloseButton,
    );

    const tagVariant = allVariants[variant];

    return (
      <div
        ref={ref}
        className={cn(
          tagVariant({ colorScheme }),
          tagSizeVariant({ size }),
          closeButtonSlot && 'pr-0',
          className,
        )}
        {...props}
      >
        {leftIconSlot}
        {labelSlot}
        {rightIconSlot}
        {closeButtonSlot && cloneElement(closeButtonSlot, { size })}
      </div>
    );
  },
);

export const TagLeftIcon = forwardRef<
  SVGElement,
  React.SVGAttributes<SVGElement>
>(({ className, children, ...rest }, ref) => {
  if (!isValidElement(children)) return <>{children}</>;

  return cloneElement(children as React.ReactElement, {
    ref,
    className: twMerge('flex items-center mr-2', className),
    ...rest,
  });
});

export const TagRightIcon = forwardRef<
  SVGElement,
  React.SVGAttributes<SVGElement>
>(({ className, children, ...rest }, ref) => {
  if (!isValidElement(children)) return <>{children}</>;

  return cloneElement(children as React.ReactElement, {
    ref,
    className: twMerge('flex items-center ml-2', className),
    ...rest,
  });
});

export const TagLabel = forwardRef<
  HTMLSpanElement,
  React.HTMLAttributes<HTMLSpanElement>
>(({ className, children, ...rest }, ref) => {
  if (!isValidElement(children))
    return <span className={twMerge(className)}>{children}</span>;

  return cloneElement(children as React.ReactElement, {
    ref,
    className: twMerge(className),
    ...rest,
  });
});

interface TagCloseButtonProps extends React.HTMLAttributes<HTMLDivElement> {
  size?: 'sm' | 'md' | 'lg';
}

export const TagCloseButton = ({
  size = 'md',
  className,
  ...props
}: TagCloseButtonProps) => {
  const iconStyle = useMemo(
    () => ({
      sm: 'size-3',
      md: 'size-4',
      lg: 'size-5',
    }),
    [size],
  )[size];

  const wrapperStyle = useMemo(
    () => ({
      sm: 'size-4',
      md: 'size-5',
      lg: 'size-6',
    }),
    [size],
  )[size];

  return (
    <span
      className={cn(
        wrapperStyle,
        'flex items-center ml-1 cursor-pointer text-grayModern-400 mr-0 bg-grayModern-100 rounded-e-md px-0.5 hover:bg-grayModern-200 hover:text-grayModern-500 transition ease-in-out',
        className,
      )}
      {...props}
    >
      <XClose className={iconStyle} />
    </span>
  );
};
