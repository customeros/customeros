import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Columns01 = ({ className, ...props }: IconProps) => (
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
      d='M6.8 3H6.2C5.0799 3 4.51984 3 4.09202 3.21799C3.71569 3.40973 3.40973 3.71569 3.21799 4.09202C3 4.51984 3 5.07989 3 6.2V17.8C3 18.9201 3 19.4802 3.21799 19.908C3.40973 20.2843 3.71569 20.5903 4.09202 20.782C4.51984 21 5.0799 21 6.2 21H6.8C7.9201 21 8.48016 21 8.90798 20.782C9.28431 20.5903 9.59027 20.2843 9.78201 19.908C10 19.4802 10 18.9201 10 17.8V6.2C10 5.0799 10 4.51984 9.78201 4.09202C9.59027 3.71569 9.28431 3.40973 8.90798 3.21799C8.48016 3 7.9201 3 6.8 3Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M17.8 3H17.2C16.0799 3 15.5198 3 15.092 3.21799C14.7157 3.40973 14.4097 3.71569 14.218 4.09202C14 4.51984 14 5.0799 14 6.2V17.8C14 18.9201 14 19.4802 14.218 19.908C14.4097 20.2843 14.7157 20.5903 15.092 20.782C15.5198 21 16.0799 21 17.2 21H17.8C18.9201 21 19.4802 21 19.908 20.782C20.2843 20.5903 20.5903 20.2843 20.782 19.908C21 19.4802 21 18.9201 21 17.8V6.2C21 5.0799 21 4.51984 20.782 4.09202C20.5903 3.71569 20.2843 3.40973 19.908 3.21799C19.4802 3 18.9201 3 17.8 3Z'
    />
  </svg>
);
