import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CreditCardRefresh = ({ className, ...props }: IconProps) => (
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
      d='M22 10H2M11 19L18.8 19C19.9201 19 20.4802 19 20.908 18.782C21.2843 18.5903 21.5903 18.2843 21.782 17.908C22 17.4802 22 16.9201 22 15.8V8.2C22 7.0799 22 6.51984 21.782 6.09202C21.5903 5.7157 21.2843 5.40974 20.908 5.21799C20.4802 5 19.9201 5 18.8 5H17M11 19L13 21M11 19L13 17M7 19H5.2C4.07989 19 3.51984 19 3.09202 18.782C2.71569 18.5903 2.40973 18.2843 2.21799 17.908C2 17.4802 2 16.9201 2 15.8V8.2C2 7.0799 2 6.51984 2.21799 6.09202C2.40973 5.71569 2.7157 5.40973 3.09202 5.21799C3.51984 5 4.0799 5 5.2 5H13M13 5L11 7M13 5L11 3'
    />
  </svg>
);
