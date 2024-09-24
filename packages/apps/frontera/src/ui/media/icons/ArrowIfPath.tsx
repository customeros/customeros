import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ArrowIfPath = ({ className, ...props }: IconProps) => (
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
      d='M3 18L6 21M6 21L9 18M6 21V10H11.75M15 18L18 21M18 21L21 18M18 21V10H11.75M11.75 10V3'
    />
  </svg>
);
