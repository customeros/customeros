import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Hurricane02 = ({ className, ...props }: IconProps) => (
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
      d='M18 12C18 15.3137 15.3137 18 12 18C8.68629 18 6 15.3137 6 12M18 12C18 8.68629 15.3137 6 12 6C8.68629 6 6 8.68629 6 12M18 12C18 16.4183 14.4183 20 10 20C5.58172 20 2 16.4183 2 12M6 12C6 7.58172 9.58172 4 14 4C18.4183 4 22 7.58172 22 12M13 12C13 12.5523 12.5523 13 12 13C11.4477 13 11 12.5523 11 12C11 11.4477 11.4477 11 12 11C12.5523 11 13 11.4477 13 12Z'
    />
  </svg>
);
