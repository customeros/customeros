import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const IntersectSquare = ({ className, ...props }: IconProps) => (
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
      d='M2 5.2C2 4.07989 2 3.51984 2.21799 3.09202C2.40973 2.71569 2.71569 2.40973 3.09202 2.21799C3.51984 2 4.0799 2 5.2 2H12.8C13.9201 2 14.4802 2 14.908 2.21799C15.2843 2.40973 15.5903 2.71569 15.782 3.09202C16 3.51984 16 4.0799 16 5.2V12.8C16 13.9201 16 14.4802 15.782 14.908C15.5903 15.2843 15.2843 15.5903 14.908 15.782C14.4802 16 13.9201 16 12.8 16H5.2C4.07989 16 3.51984 16 3.09202 15.782C2.71569 15.5903 2.40973 15.2843 2.21799 14.908C2 14.4802 2 13.9201 2 12.8V5.2Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M8 11.2C8 10.0799 8 9.51984 8.21799 9.09202C8.40973 8.71569 8.71569 8.40973 9.09202 8.21799C9.51984 8 10.0799 8 11.2 8H18.8C19.9201 8 20.4802 8 20.908 8.21799C21.2843 8.40973 21.5903 8.71569 21.782 9.09202C22 9.51984 22 10.0799 22 11.2V18.8C22 19.9201 22 20.4802 21.782 20.908C21.5903 21.2843 21.2843 21.5903 20.908 21.782C20.4802 22 19.9201 22 18.8 22H11.2C10.0799 22 9.51984 22 9.09202 21.782C8.71569 21.5903 8.40973 21.2843 8.21799 20.908C8 20.4802 8 19.9201 8 18.8V11.2Z'
    />
  </svg>
);
