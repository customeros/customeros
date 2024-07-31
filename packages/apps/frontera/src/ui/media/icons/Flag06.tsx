import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Flag06 = ({ className, ...props }: IconProps) => (
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
      d='M4 21L4 4M4 13H11.4C11.9601 13 12.2401 13 12.454 12.891C12.6422 12.7951 12.7951 12.6422 12.891 12.454C13 12.2401 13 11.9601 13 11.4V4.6C13 4.03995 13 3.75992 12.891 3.54601C12.7951 3.35785 12.6422 3.20487 12.454 3.10899C12.2401 3 11.9601 3 11.4 3H5.6C5.03995 3 4.75992 3 4.54601 3.10899C4.35785 3.20487 4.20487 3.35785 4.10899 3.54601C4 3.75992 4 4.03995 4 4.6V13ZM13 5H19.4C19.9601 5 20.2401 5 20.454 5.10899C20.6422 5.20487 20.7951 5.35785 20.891 5.54601C21 5.75992 21 6.03995 21 6.6V13.4C21 13.9601 21 14.2401 20.891 14.454C20.7951 14.6422 20.6422 14.7951 20.454 14.891C20.2401 15 19.9601 15 19.4 15H14.6C14.0399 15 13.7599 15 13.546 14.891C13.3578 14.7951 13.2049 14.6422 13.109 14.454C13 14.2401 13 13.9601 13 13.4V5Z'
    />
  </svg>
);
