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
  | 'copy-03'
  | 'download-04'
  | 'dot-single'
  | 'dots-vertical'
  | 'dot-live-success'
  | 'dot-live-primary'
  | 'cloud-off'
  | 'invoice'
  | 'invoice-upcoming'
  | 'invoice-check'
  | 'building-07'
  | 'check-heart'
  | 'radio-dot'
  | 'users-01'
  | 'heart-hand'
  | 'radar'
  | 'signature'
  | 'target-05'
  | 'plus-circle'
  | 'coins-stacked-01'
  | 'shuffle-01'
  | 'play'
  | 'log-out-01'
  | 'layers-two-01'
  | 'settings-02'
  | 'x-close'
  | 'share-07'
  | 'tag-01'
  | 'activity-heart'
  | 'seeding'
  | 'message-x-circle'
  | 'broken-heart'
  | 'align-horizontal-centre-02'
  | 'check-verified-02'
  | 'target-04'
  | 'users-02'
  | 'key-01'
  | 'eye-off'
  | 'eye'
  | 'columns-03'
  | 'user-03'
  | 'user-01'
  | 'activity'
  | 'thumbs-down'
  | 'message-question-circle'
  | 'x-circle';

interface IconProps extends SVGAttributes<SVGElement> {
  name: IconName;
  className?: string;
}

export const Icon = ({
  name,
  fill,
  width,
  height,
  stroke,
  className,
  strokeWidth,
  ...props
}: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    width={width ?? 24}
    fill={fill ?? 'none'}
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
