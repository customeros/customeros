import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CornerLeftUp = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      d='M21 20H17.6C14.2397 20 12.5595 20 11.2761 19.346C10.1471 18.7708 9.2292 17.8529 8.65396 16.7239C8 15.4405 8 13.7603 8 10.4V4M8 4L13 9M8 4L3 9'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
