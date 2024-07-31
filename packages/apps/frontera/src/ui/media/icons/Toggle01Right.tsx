import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Toggle01Right = ({ className, ...props }: IconProps) => (
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
      d='M17 17H7C4.23858 17 2 14.7614 2 12C2 9.23858 4.23858 7 7 7H17M17 17C19.7614 17 22 14.7614 22 12C22 9.23858 19.7614 7 17 7M17 17C14.2386 17 12 14.7614 12 12C12 9.23858 14.2386 7 17 7'
    />
  </svg>
);
