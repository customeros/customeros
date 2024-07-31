import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Microscope = ({ className, ...props }: IconProps) => (
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
      d='M3 22H12M11 6.25204C11.6392 6.08751 12.3094 6 13 6C17.4183 6 21 9.58172 21 14C21 17.3574 18.9318 20.2317 16 21.4185M5.5 13H9.5C9.96466 13 10.197 13 10.3902 13.0384C11.1836 13.1962 11.8038 13.8164 11.9616 14.6098C12 14.803 12 15.0353 12 15.5C12 15.9647 12 16.197 11.9616 16.3902C11.8038 17.1836 11.1836 17.8038 10.3902 17.9616C10.197 18 9.96466 18 9.5 18H5.5C5.03534 18 4.80302 18 4.60982 17.9616C3.81644 17.8038 3.19624 17.1836 3.03843 16.3902C3 16.197 3 15.9647 3 15.5C3 15.0353 3 14.803 3.03843 14.6098C3.19624 13.8164 3.81644 13.1962 4.60982 13.0384C4.80302 13 5.03534 13 5.5 13ZM4 5.5V13H11V5.5C11 3.567 9.433 2 7.5 2C5.567 2 4 3.567 4 5.5Z'
    />
  </svg>
);
