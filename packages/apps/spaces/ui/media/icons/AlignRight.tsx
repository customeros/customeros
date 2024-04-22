import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const AlignRight = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      d='M21 10H8M21 6H4M21 14H4M21 18H8'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
