import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Wind03 = ({ className, ...props }: IconProps) => (
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
      d='M16.7639 6.5C17.3132 5.88625 18.1115 5.5 19 5.5C20.6569 5.5 22 6.84315 22 8.5C22 10.1569 20.6569 11.5 19 11.5H13M6.7639 4C7.31322 3.38625 8.1115 3 9 3C10.6569 3 12 4.34315 12 6C12 7.65685 10.6569 9 9 9H2M10.7639 20C11.3132 20.6137 12.1115 21 13 21C14.6569 21 16 19.6569 16 18C16 16.3431 14.6569 15 13 15H2'
    />
  </svg>
);
