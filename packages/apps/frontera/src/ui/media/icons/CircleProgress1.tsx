import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CircleProgress1 = ({ className, ...props }: IconProps) => (
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
      d='M11.9999 2C15.0413 2 17.7655 3.35767 19.5995 5.5C21.0961 7.24803 22 9.51846 22 12C22 17.5228 17.5228 22 11.9999 22C8.50892 22 5.43571 20.2111 3.64697 17.5'
    />
  </svg>
);
