import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Hurricane01 = ({ className, ...props }: IconProps) => (
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
      d='M16.5 20.5002C15.2465 20.814 13.6884 21 12 21C10.3116 21 8.75349 20.814 7.5 20.5002M18 16.4305C16.5341 16.9842 14.3894 17.3333 12 17.3333C9.61061 17.3333 7.46589 16.9842 6 16.4305M4.5 11.6679C5.93143 12.5598 8.75311 13.1667 12 13.1667C15.2469 13.1667 18.0686 12.5598 19.5 11.6679M21 6C21 7.65685 16.9706 9 12 9C7.02944 9 3 7.65685 3 6C3 4.34315 7.02944 3 12 3C16.9706 3 21 4.34315 21 6Z'
    />
  </svg>
);
