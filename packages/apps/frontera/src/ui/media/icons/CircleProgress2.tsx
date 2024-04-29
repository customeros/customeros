import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CircleProgress2 = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <circle cx='12' cy='12' r='10' stroke='#D0D5DD' strokeWidth='2' />
    <path
      d='M11.9999 2C15.0413 2 17.7655 3.35767 19.5995 5.5C21.0961 7.24803 22 9.51846 22 12C22 17.5228 17.5228 22 11.9999 22C8.50892 22 5.43571 20.2111 3.64697 17.5'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
