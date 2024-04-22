import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ChevronRightDouble = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      d='M6 17L11 12L6 7M13 17L18 12L13 7'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
