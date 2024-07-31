import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CloudSun03 = ({ className, ...props }: IconProps) => (
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
      d='M3.15003 11C3.05165 10.5153 3 10.0137 3 9.5C3 5.35786 6.35786 2 10.5 2C14.3031 2 17.445 4.83064 17.9339 8.5M6 22C3.79086 22 2 20.2091 2 18C2 15.7909 3.79086 14 6 14C6.11333 14 6.22556 14.0047 6.3365 14.014C7.15622 11.6763 9.38235 10 12 10C14.2248 10 16.1668 11.2109 17.2029 13.0097C17.3011 13.0033 17.4002 13 17.5 13C19.9853 13 22 15.0147 22 17.5C22 19.9853 19.9853 22 17.5 22C13.7633 22 10.0546 22 6 22Z'
    />
  </svg>
);
