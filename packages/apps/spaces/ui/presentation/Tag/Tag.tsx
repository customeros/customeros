import { cloneElement, isValidElement } from 'react';

import { twMerge } from 'tailwind-merge';
import { VariantProps } from 'class-variance-authority';

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

export const Tag = ({
  size,
  children,
  className,
  colorScheme,
  variant = 'outline',
  ...props
}: TagProps) => {
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
      className={twMerge(tagVariant({ colorScheme }), tagSizeVariant({ size }))}
      {...props}
    >
      {leftIconSlot}
      {labelSlot}
      {rightIconSlot}
      {closeButtonSlot}
    </div>
  );
};

export const TagLeftIcon = ({ className, children, ...rest }: TagProps) => {
  if (!isValidElement(children)) return <>{children}</>;

  return cloneElement(children as React.ReactElement, {
    className: twMerge('flex items-center mr-2', className),
    ...rest,
  });
};

export const TagRightIcon = ({ className, children, ...rest }: TagProps) => {
  if (!isValidElement(children)) return <>{children}</>;

  return cloneElement(children as React.ReactElement, {
    className: twMerge('flex items-center ml-2', className),
    ...rest,
  });
};

export const TagLabel = ({ className, children, ...rest }: TagProps) => {
  if (!isValidElement(children))
    return <span className={twMerge(className)}>{children}</span>;

  return cloneElement(children as React.ReactElement, {
    className: twMerge(className),
    ...rest,
  });
};

export const TagCloseButton = (props: React.HTMLAttributes<SVGAElement>) => {
  return <XClose {...props} />;
};
