import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const BellOff01 = ({ className, ...props }: IconProps) => (
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
      d='M8.63306 3.03371C9.61959 2.3649 10.791 2 12 2C13.5913 2 15.1174 2.63214 16.2426 3.75736C17.3679 4.88258 18 6.4087 18 8C18 10.1008 18.2702 11.7512 18.6484 13.0324M6.25867 6.25724C6.08866 6.81726 6 7.40406 6 8C6 11.0902 5.22047 13.206 4.34966 14.6054C3.61513 15.7859 3.24786 16.3761 3.26132 16.5408C3.27624 16.7231 3.31486 16.7926 3.46178 16.9016C3.59446 17 4.19259 17 5.38885 17H17M9.35418 21C10.0593 21.6224 10.9856 22 12 22C13.0144 22 13.9407 21.6224 14.6458 21M21 21L3 3'
    />
  </svg>
);
