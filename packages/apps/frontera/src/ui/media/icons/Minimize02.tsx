import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Minimize02 = ({ className, ...props }: IconProps) => (
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
      d='M2.99988 8H3.19988C4.88004 8 5.72011 8 6.36185 7.67302C6.92634 7.3854 7.38528 6.92646 7.6729 6.36197C7.99988 5.72024 7.99988 4.88016 7.99988 3.2V3M2.99988 16H3.19988C4.88004 16 5.72011 16 6.36185 16.327C6.92634 16.6146 7.38528 17.0735 7.6729 17.638C7.99988 18.2798 7.99988 19.1198 7.99988 20.8V21M15.9999 3V3.2C15.9999 4.88016 15.9999 5.72024 16.3269 6.36197C16.6145 6.92646 17.0734 7.3854 17.6379 7.67302C18.2796 8 19.1197 8 20.7999 8H20.9999M15.9999 21V20.8C15.9999 19.1198 15.9999 18.2798 16.3269 17.638C16.6145 17.0735 17.0734 16.6146 17.6379 16.327C18.2796 16 19.1197 16 20.7999 16H20.9999'
    />
  </svg>
);
