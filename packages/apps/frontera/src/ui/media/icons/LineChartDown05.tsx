import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const LineChartDown05 = ({ className, ...props }: IconProps) => (
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
      d='M6 10L9.43431 13.4343C9.63232 13.6323 9.73133 13.7313 9.84549 13.7684C9.94591 13.8011 10.0541 13.8011 10.1545 13.7684C10.2687 13.7313 10.3677 13.6323 10.5657 13.4343L13.4343 10.5657C13.6323 10.3677 13.7313 10.2687 13.8455 10.2316C13.9459 10.1989 14.0541 10.1989 14.1545 10.2316C14.2687 10.2687 14.3677 10.3677 14.5657 10.5657L18 14M22 12C22 17.5228 17.5228 22 12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12Z'
    />
  </svg>
);
