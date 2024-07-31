import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Subscript = ({ className, ...props }: IconProps) => (
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
      d='M3 5L13 15M13 5L3 15M21 19.0001H17C17 17.5001 17.44 17.0001 18.5 16.5001C19.56 16.0001 21 15.3301 21 14.0001C21 13.5301 20.83 13.0701 20.52 12.7101C20.1999 12.3476 19.7649 12.1061 19.2879 12.026C18.811 11.9459 18.321 12.0321 17.9 12.2701C17.48 12.5101 17.16 12.8901 17 13.3401'
    />
  </svg>
);
