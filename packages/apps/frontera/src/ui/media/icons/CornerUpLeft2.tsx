import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CornerUpLeft2 = ({ className, ...props }: IconProps) => (
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
      d='M11 14L6 9M6 9L11 4M6 9H12.4C15.7603 9 17.4405 9 18.7239 9.65396C19.8529 10.2292 20.7708 11.1471 21.346 12.2761C22 13.5595 22 15.2397 22 18.6V20M7 14L2 9L7 4'
    />
  </svg>
);
