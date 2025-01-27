import type { SVGAttributes } from 'react';

import { twMerge } from 'tailwind-merge';

export type IconName =
  | 'atom-01'
  | 'alert-circle'
  | 'archive'
  | 'arrow-dropdown'
  | 'bubbles'
  | 'chevron-down'
  | 'chevron-left'
  | 'chevron-right'
  | 'download-04'
  | 'dots-vertical'
  | 'cloud-off'
  | 'invoice'
  | 'invoice-upcoming'
  | 'invoice-check'
  | 'building-07'
  | 'check-heart'
  | 'users-01'
  | 'heart-hand'
  | 'signature'
  | 'target-05'
  | 'coins-stacked-01'
  | 'shuffle-01'
  | 'play'
  | 'log-out-01'
  | 'layers-two-01'
  | 'settings-02'
  | 'x-close';

interface IconProps extends SVGAttributes<SVGElement> {
  name: IconName;
  className?: string;
}

export const Icon = ({
  name,
  width,
  height,
  stroke,
  className,
  strokeWidth,
  ...props
}: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    width={width ?? 24}
    height={height ?? 24}
    strokeLinecap='round'
    strokeLinejoin='round'
    strokeWidth={strokeWidth ?? 2}
    stroke={stroke ?? 'currentColor'}
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <use xlinkHref={`/icons.svg#${name}`} />
  </svg>
);
