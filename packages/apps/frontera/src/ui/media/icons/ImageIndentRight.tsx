import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ImageIndentRight = ({ className, ...props }: IconProps) => (
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
      d='M21 4H3M21 20H3M9 9.25H3M9 14.75H3M14.6 16H19.4C19.9601 16 20.2401 16 20.454 15.891C20.6422 15.7951 20.7951 15.6422 20.891 15.454C21 15.2401 21 14.9601 21 14.4V9.6C21 9.03995 21 8.75992 20.891 8.54601C20.7951 8.35785 20.6422 8.20487 20.454 8.10899C20.2401 8 19.9601 8 19.4 8H14.6C14.0399 8 13.7599 8 13.546 8.10899C13.3578 8.20487 13.2049 8.35785 13.109 8.54601C13 8.75992 13 9.03995 13 9.6V14.4C13 14.9601 13 15.2401 13.109 15.454C13.2049 15.6422 13.3578 15.7951 13.546 15.891C13.7599 16 14.0399 16 14.6 16Z'
    />
  </svg>
);
