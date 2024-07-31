import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const LineChartUp05 = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M18 10L14.5657 13.4343C14.3677 13.6323 14.2687 13.7313 14.1545 13.7684C14.0541 13.8011 13.9459 13.8011 13.8455 13.7684C13.7313 13.7313 13.6323 13.6323 13.4343 13.4343L10.5657 10.5657C10.3677 10.3677 10.2687 10.2687 10.1545 10.2316C10.0541 10.1989 9.94591 10.1989 9.84549 10.2316C9.73133 10.2687 9.63232 10.3677 9.43431 10.5657L6 14M22 12C22 17.5228 17.5228 22 12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12Z'
    />
  </svg>
);
