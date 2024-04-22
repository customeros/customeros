import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CornerDownRight = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      d='M4 4V5.4C4 8.76031 4 10.4405 4.65396 11.7239C5.2292 12.8529 6.14708 13.7708 7.27606 14.346C8.55953 15 10.2397 15 13.6 15H20M20 15L15 10M20 15L15 20'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
