import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const RefreshCcw03 = ({ className, ...props }: IconProps) => (
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
      d='M14 2C14 2 14.8492 2.12132 18.364 5.63604C21.8787 9.15076 21.8787 14.8492 18.364 18.364C17.1187 19.6092 15.5993 20.4133 14 20.7762M14 2L20 2M14 2L14 8M10 21.9998C10 21.9998 9.15076 21.8785 5.63604 18.3638C2.12132 14.849 2.12132 9.15056 5.63604 5.63584C6.88131 4.39057 8.40072 3.5865 10 3.22363M10 21.9998L4 22M10 21.9998L10 16'
    />
  </svg>
);
